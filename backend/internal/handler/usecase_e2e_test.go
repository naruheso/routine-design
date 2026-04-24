package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/amical/routine-design/backend/internal/middleware"
	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
	"github.com/amical/routine-design/backend/internal/service"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	pool, err := repository.NewDBPool(ctx)
	if err != nil {
		fmt.Printf("E2Eテスト用DB接続失敗: %v\n", err)
		os.Exit(1)
	}
	testPool = pool
	code := m.Run()
	testPool.Close()
	os.Exit(code)
}

func setupE2E(t *testing.T) *gin.Engine {
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

	authService := service.NewAuthService(userRepo)
	expenseService := service.NewExpenseService(expenseRepo, approvalRepo)
	approvalService := service.NewApprovalService(expenseRepo, approvalRepo)
	reportService := service.NewReportService(reportRepo)

	authHandler := NewAuthHandler(authService)
	expenseHandler := NewExpenseHandler(expenseService)
	approvalHandler := NewApprovalHandler(approvalService)
	categoryHandler := NewCategoryHandler(categoryRepo)
	reportHandler := NewReportHandler(reportService)

	// 未使用エラー回避
	_ = categoryHandler
	_ = reportHandler

	r := gin.New()
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			expenses := authorized.Group("/expenses")
			{
				expenses.POST("", expenseHandler.Create)
				expenses.POST("/:id/submit", expenseHandler.Submit)
				expenses.GET("/:id", expenseHandler.GetByID)
			}
			approvals := authorized.Group("/approvals")
			{
				approvals.POST("/:id/approve", approvalHandler.Approve)
			}
		}
	}
	return r
}

func TestUseCase_ExpenseFlow_E2E(t *testing.T) {
	r := setupE2E(t)
	ctx := context.Background()

	// 1. マスターデータ準備
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	passHash := string(hashed)
	
	userUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", userUID, "一般ユーザー", "user@test.com", passHash)
	
	managerUID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)", managerUID, "マネージャー", "mgr@test.com", passHash)
	testPool.Exec(ctx, "INSERT INTO user_roles (user_id, role) VALUES ($1, $2)", managerUID, "manager")

	catID := uuid.New()
	testPool.Exec(ctx, "INSERT INTO categories (id, code, name) VALUES ($1, $2, $3)", catID, "test-cat", "テスト科目")

	// --- Step 1: 一般ユーザーでログイン ---
	loginReq := model.LoginRequest{Email: "user@test.com", Password: "password"}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	var loginRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &loginRes)
	userToken := loginRes.Data.Token
	require.NotEmpty(t, userToken)

	// --- Step 2: 経費作成 (下書き) ---
	createReq := model.CreateExpenseRequest{
		ExpenseDate: "2026-04-24",
		CategoryID:  catID.String(),
		Amount:      1200,
		Description: "E2Eテスト用経費",
	}
	body, _ = json.Marshal(createReq)
	req, _ = http.NewRequest("POST", "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var createRes struct {
		Data struct { ID string `json:"id"` } `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &createRes)
	expenseID := createRes.Data.ID
	require.NotEmpty(t, expenseID)

	// --- Step 3: 申請 ---
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/expenses/%s/submit", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// --- Step 4: マネージャーでログイン ---
	loginReq = model.LoginRequest{Email: "mgr@test.com", Password: "password"}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &loginRes)
	mgrToken := loginRes.Data.Token

	// --- Step 5: 承認 (1次承認: マネージャー) ---
	comment := "承認します"
	approveReq := model.ApprovalRequest{Comment: &comment}
	body, _ = json.Marshal(approveReq)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/approvals/%s/approve", expenseID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+mgrToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// --- Step 6: ステータス確認 (pending_expense_admin になっているか) ---
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	var finalRes struct {
		Data struct {
			Expense struct {
				Status string `json:"status"`
			} `json:"expense"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &finalRes)
	assert.Equal(t, "pending_expense_admin", finalRes.Data.Expense.Status)
}
