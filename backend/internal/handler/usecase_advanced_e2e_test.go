package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
	"github.com/amical/routine-design/backend/internal/service"
)

// setupFullE2E は全ユースケースをカバーするために全エンドポイントを登録したルーターを返す。
func setupFullE2E(t *testing.T) *gin.Engine {
	ctx := context.Background()
	// クリーンアップ
	tables := []string{"approval_histories", "expenses", "categories", "user_roles", "users"}
	for _, table := range tables {
		testPool.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
	}

	// 各レイヤー初期化
	userRepo := repository.NewUserRepository(testPool)
	expenseRepo := repository.NewExpenseRepository(testPool)
	approvalRepo := repository.NewApprovalRepository(testPool)
	categoryRepo := repository.NewCategoryRepository(testPool)
	reportRepo := repository.NewReportRepository(testPool)
	txManager := repository.NewTransactionManager(testPool)

	authService := service.NewAuthService(userRepo)
	expenseService := service.NewExpenseService(expenseRepo, approvalRepo, txManager)
	approvalService := service.NewApprovalService(expenseRepo, approvalRepo, txManager)
	reportService := service.NewReportService(reportRepo)

	authHandler := NewAuthHandler(authService)
	expenseHandler := NewExpenseHandler(expenseService)
	approvalHandler := NewApprovalHandler(approvalService)
	categoryHandler := NewCategoryHandler(categoryRepo)
	reportHandler := NewReportHandler(reportService)

	r := SetupRouter(authHandler, expenseHandler, approvalHandler, categoryHandler, reportHandler)
	return r
}

func TestUseCase_FullLifecycle_E2E(t *testing.T) {
	r := setupFullE2E(t)
	ctx := context.Background()

	// 1. マスターデータ準備
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	
	userUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", userUID, "山田太郎", "yamada@test.com", string(passHash))
	
	managerUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", managerUID, "マネージャー", "mgr@test.com", string(passHash))
	testPool.Exec(ctx, "INSERT INTO user_roles (user_id, role) VALUES ($1, $2)", managerUID, "manager")

	adminUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", adminUID, "経理部長", "admin@test.com", string(passHash))
	// 全ロールを付与して一気に承認可能にする
	testPool.Exec(ctx, "INSERT INTO user_roles (user_id, role) VALUES ($1, 'manager'), ($1, 'expense_admin'), ($1, 'finance_director')", adminUID)

	catID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO categories (id, code, name) VALUES ($1, $2, $3)", catID, "trans", "交通費")

	// トークン取得用ヘルパー
	getToken := func(email string) string {
		loginReq := model.LoginRequest{Email: email, Password: "password"}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var res struct { Data struct { Token string `json:"token"` } `json:"data"` }
		json.Unmarshal(w.Body.Bytes(), &res)
		return res.Data.Token
	}

	userToken := getToken("yamada@test.com")
	mgrToken := getToken("mgr@test.com")
	adminToken := getToken("admin@test.com")

	// --- Step 1: 経費作成 & 申請 (UC-03) ---
	createReq := model.CreateExpenseRequest{
		ExpenseDate: "2026-04-24", CategoryID: catID.String(), Amount: 5000, Description: "出張費",
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var createRes struct { Data struct { ID string `json:"id"` } `json:"data"` }
	json.Unmarshal(w.Body.Bytes(), &createRes)
	expenseID := createRes.Data.ID

	// 申請
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/expenses/%s/submit", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(httptest.NewRecorder(), req)

	// --- Step 2: マネージャーが差し戻す (UC-07, UC-09) ---
	// 差し戻し (UC-09)
	returnMsg := "領収書が不足しています"
	returnReq := model.ApprovalRequest{Comment: &returnMsg}
	body, _ = json.Marshal(returnReq)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/approvals/%s/return", expenseID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+mgrToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// --- Step 3: ユーザーが修正・再申請 (UC-05, UC-04) ---
	// 修正 (UC-04)
	updateReq := model.UpdateExpenseRequest{
		ExpenseDate: "2026-04-24", CategoryID: catID.String(), Amount: 4500, Description: "出張費(金額修正)",
	}
	body, _ = json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/expenses/%s", expenseID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(httptest.NewRecorder(), req)

	// 再申請
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/expenses/%s/submit", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(httptest.NewRecorder(), req)

	// --- Step 4: 最終承認まで進める (UC-08) ---
	// adminは全ロールを持っているので、3回承認すれば approved になるはず
	for i := 0; i < 3; i++ {
		req, _ = http.NewRequest("POST", fmt.Sprintf("/api/approvals/%s/approve", expenseID), nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	// ステータス確認
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "approved")

	// --- Step 5: 集計レポートに反映されているか (UC-11) ---
	req, _ = http.NewRequest("GET", "/api/reports/by-employee?year=2026&month=4", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	// 修正後の金額 4500 が集計に含まれているか
	assert.Contains(t, w.Body.String(), "4500")
	assert.Contains(t, w.Body.String(), "山田太郎")
}

func TestUseCase_RejectionAndReports_E2E(t *testing.T) {
	r := setupFullE2E(t)
	ctx := context.Background()

	// 1. マスターデータ準備
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	userUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", userUID, "田中次郎", "tanaka@test.com", string(passHash))
	
	managerUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", managerUID, "マネージャー", "mgr@test.com", string(passHash))
	testPool.Exec(ctx, "INSERT INTO user_roles (user_id, role) VALUES ($1, $2)", managerUID, "manager")

	catID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO categories (id, code, name) VALUES ($1, $2, $3)", catID, "supplies", "事務用品費")

	getToken := func(email string) string {
		loginReq := model.LoginRequest{Email: email, Password: "password"}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var res struct { Data struct { Token string `json:"token"` } `json:"data"` }
		json.Unmarshal(w.Body.Bytes(), &res)
		return res.Data.Token
	}
	userToken := getToken("tanaka@test.com")
	mgrToken := getToken("mgr@test.com")

	// --- Step 1: 申請 ---
	createReq := model.CreateExpenseRequest{
		ExpenseDate: "2026-04-25", CategoryID: catID.String(), Amount: 3000, Description: "マウス購入",
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var createRes struct { Data struct { ID string `json:"id"` } `json:"data"` }
	json.Unmarshal(w.Body.Bytes(), &createRes)
	expenseID := createRes.Data.ID

	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/expenses/%s/submit", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(httptest.NewRecorder(), req)

	// --- Step 2: マネージャーが否認する (UC-10) ---
	rejectMsg := "マウスは会社から支給されるものを利用してください"
	rejectReq := model.ApprovalRequest{Comment: &rejectMsg}
	body, _ = json.Marshal(rejectReq)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/approvals/%s/reject", expenseID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+mgrToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// ステータス確認
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "rejected")

	// --- Step 3: BR-01検証 (否認済みは更新不可) ---
	updateReq := model.UpdateExpenseRequest{
		ExpenseDate: "2026-04-25", CategoryID: catID.String(), Amount: 1000, Description: "キーボードに修正",
	}
	body, _ = json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/expenses/%s", expenseID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// 400 Bad Request または 403 Forbidden が返るはず (ビジネスルール違反)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "編集できないステータス")

	// --- Step 4: カテゴリ別集計レポートの確認 (UC-11) ---
	req, _ = http.NewRequest("GET", "/api/reports/by-category?year=2026&month=4", nil)
	req.Header.Set("Authorization", "Bearer "+userToken) // 閲覧権限がある前提
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// 否認されたマウス（3000円）は集計に含まれていないはず（totalCount/grandTotalが0）
	assert.Contains(t, w.Body.String(), "\"totalCount\":0")
}
