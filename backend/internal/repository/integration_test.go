package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

// TestMain はテスト全体の実行前後で一度だけ呼ばれる。
func TestMain(m *testing.M) {
	// DB接続の確立
	ctx := context.Background()
	pool, err := NewDBPool(ctx)
	if err != nil {
		fmt.Printf("テスト用DB接続失敗: %v\n", err)
		os.Exit(1)
	}
	testPool = pool

	// テスト実行
	code := m.Run()

	// 終了処理
	testPool.Close()
	os.Exit(code)
}

// setupTestData は各テストの前にテーブルを空にする。
func setupTestData(ctx context.Context, t *testing.T) {
	t.Helper()
	// 参照整合性の順序を考慮して削除
	tables := []string{"approval_histories", "expenses", "categories", "users"}
	for _, table := range tables {
		if _, err := testPool.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			t.Fatalf("table %s のクリーンアップに失敗: %v", table, err)
		}
	}
}
