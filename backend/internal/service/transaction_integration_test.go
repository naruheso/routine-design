package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
	"github.com/amical/routine-design/backend/internal/service"
)

func TestExpenseService_Submit_TransactionRollback(t *testing.T) {
	ctx := context.Background()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "routine_design"),
	)
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// 1. マスターデータ準備
	userUID := uuid.New().String()
	catID := uuid.New().String()
	_, err = pool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash) VALUES ($1, 'Txテスト社員', 'tx@example.com', 'hash')", userUID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, "INSERT INTO categories (id, code, name) VALUES ($1, 'TX', 'TXテスト科目')", catID)
	require.NoError(t, err)

	expenseRepo := repository.NewExpenseRepository(pool)
	approvalRepo := repository.NewApprovalRepository(pool)
	
	// TransactionManagerの生成
	txManager := repository.NewTransactionManager(pool)

	// ここでは、サービスに直接TransactionManagerを渡せるようにする前提でテストを書く。
	// (現状のService層にはまだ実装していないので、テストをパスさせるための改修が必要になる)
	expenseService := service.NewExpenseService(expenseRepo, approvalRepo, txManager)

	// 経費下書き作成
	createReq := model.CreateExpenseRequest{
		ExpenseDate: "2026-04-24", CategoryID: catID, Amount: 1000, Description: "Txテスト",
	}
	expense, err := expenseService.Create(ctx, userUID, createReq)
	require.NoError(t, err)
	require.Equal(t, model.StatusDraft, expense.Status)

	defer func() {
		pool.Exec(ctx, "DELETE FROM expenses WHERE user_id = $1", userUID)
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userUID)
		pool.Exec(ctx, "DELETE FROM categories WHERE id = $1", catID)
	}()

	// 2. 意図的なエラーの誘発（テーブルを一時的にリネーム）
	_, err = pool.Exec(ctx, "ALTER TABLE approval_histories RENAME TO approval_histories_temp")
	require.NoError(t, err)

	defer func() {
		// テスト終了時に必ず元に戻す
		pool.Exec(ctx, "ALTER TABLE approval_histories_temp RENAME TO approval_histories")
	}()

	// 3. 申請処理を実行（ステータス更新は成功するが、履歴作成でDBエラーになるはず）
	_, err = expenseService.Submit(ctx, expense.ID, userUID)
	
	// エラーが発生していることを確認
	require.Error(t, err)
	assert.Contains(t, err.Error(), "relation \"approval_histories\" does not exist")

	// 4. ロールバックの検証（ステータスが draft のままであること）
	// 別コネクション(新しいコンテキスト)で確認
	checkCtx := context.Background()
	savedExpense, err := expenseRepo.FindByID(checkCtx, expense.ID)
	require.NoError(t, err)
	
	// もしトランザクションが効いていなければ "pending_manager" になってしまう
	assert.Equal(t, model.StatusDraft, savedExpense.Status, "トランザクションがロールバックされていません！")
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
