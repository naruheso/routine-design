package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/amical/routine-design/backend/internal/model"
)

func TestReportRepository_Aggregate_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	setupTestData(ctx, t)
	repo := NewReportRepository(testPool)

	// 1. テストデータ準備 (Seed)
	// ユーザー作成
	user1ID := uuid.New()
	_, err := testPool.Exec(ctx, "INSERT INTO users (id, email, password_hash, name) VALUES ($1, $2, $3, $4)",
		user1ID, "user1@example.com", "hash", "山田太郎")
	require.NoError(t, err)

	user2ID := uuid.New()
	_, err = testPool.Exec(ctx, "INSERT INTO users (id, email, password_hash, name) VALUES ($1, $2, $3, $4)",
		user2ID, "user2@example.com", "hash", "佐藤花子")
	require.NoError(t, err)

	// カテゴリ作成
	cat1ID := uuid.New()
	_, err = testPool.Exec(ctx, "INSERT INTO categories (id, code, name, sort_order) VALUES ($1, $2, $3, $4)",
		cat1ID, "trans-test", "交通費", 1)
	require.NoError(t, err)

	cat2ID := uuid.New()
	_, err = testPool.Exec(ctx, "INSERT INTO categories (id, code, name, sort_order) VALUES ($1, $2, $3, $4)",
		cat2ID, "food-test", "接待交際費", 2)
	require.NoError(t, err)

	// 経費作成 (集計対象: 2026年4月)
	targetDate := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)
	otherDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		uid      uuid.UUID
		cid      uuid.UUID
		amount   int
		date     time.Time
		status   string
		isTarget bool
	}{
		{user1ID, cat1ID, 1000, targetDate, "approved", true},  // 対象
		{user1ID, cat1ID, 2000, targetDate, "approved", true},  // 対象
		{user2ID, cat2ID, 5000, targetDate, "approved", true},  // 対象
		{user1ID, cat1ID, 9999, targetDate, "draft", false},    // 除外 (status)
		{user1ID, cat1ID, 9999, otherDate, "approved", false},  // 除外 (month)
	}

	for _, tc := range testCases {
		_, err := testPool.Exec(ctx, `
			INSERT INTO expenses (id, user_id, category_id, amount, expense_date, status, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.New(), tc.uid, tc.cid, tc.amount, tc.date, tc.status, "test")
		require.NoError(t, err)
	}

	// 2. 実行: 社員別集計
	t.Run("社員別集計の検証", func(t *testing.T) {
		reports, err := repo.AggregateByEmployee(ctx, 2026, 4)
		require.NoError(t, err)

		// 山田太郎: exp-1(1000) + exp-2(2000) = 3000
		// 佐藤花子: exp-3(5000) = 5000
		assert.Len(t, reports, 2)
		
		var user1Report, user2Report *model.EmployeeReport
		for i := range reports {
			if reports[i].UserID == user1ID.String() { user1Report = &reports[i] }
			if reports[i].UserID == user2ID.String() { user2Report = &reports[i] }
		}

		require.NotNil(t, user1Report)
		assert.Equal(t, 2, user1Report.Count)
		assert.Equal(t, 3000, user1Report.TotalAmount)

		require.NotNil(t, user2Report)
		assert.Equal(t, 1, user2Report.Count)
		assert.Equal(t, 5000, user2Report.TotalAmount)
	})

	// 3. 実行: 勘定科目別集計
	t.Run("勘定科目別集計の検証", func(t *testing.T) {
		reports, err := repo.AggregateByCategory(ctx, 2026, 4)
		require.NoError(t, err)

		// 交通費(cat-1): exp-1(1000) + exp-2(2000) = 3000
		// 接待交際費(cat-2): exp-3(5000) = 5000
		assert.Len(t, reports, 2)

		var cat1Report, cat2Report *model.CategoryReport
		for i := range reports {
			if reports[i].CategoryID == cat1ID.String() { cat1Report = &reports[i] }
			if reports[i].CategoryID == cat2ID.String() { cat2Report = &reports[i] }
		}

		require.NotNil(t, cat1Report)
		assert.Equal(t, 3000, cat1Report.TotalAmount)
		assert.Equal(t, 2, cat1Report.Count)

		require.NotNil(t, cat2Report)
		assert.Equal(t, 5000, cat2Report.TotalAmount)
	})
}
