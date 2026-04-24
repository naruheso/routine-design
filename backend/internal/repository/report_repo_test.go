package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportRepository_AggregateByEmployee(t *testing.T) {
	t.Parallel()

	// pgxmockのプールを作成
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewReportRepository(mock)
	ctx := context.Background()

	tests := []struct {
		name       string
		year       int
		month      int
		mockRows   *pgxmock.Rows
		wantErr    bool
		wantCount  int
		wantTotal  int
	}{
		{
			name:  "正常系: 2件のデータが返る",
			year:  2026,
			month: 4,
			mockRows: pgxmock.NewRows([]string{"user_id", "user_name", "count", "total_amount"}).
				AddRow("u1", "山田太郎", 2, 3000).
				AddRow("u2", "佐藤花子", 1, 5000),
			wantCount: 2,
			wantTotal: 8000,
		},
		{
			name:      "正常系: 0件",
			year:      2026,
			month:     5,
			mockRows:  pgxmock.NewRows([]string{"user_id", "user_name", "count", "total_amount"}),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SQLの期待値を設定
			// 正規表現でクエリをマッチさせる
			mock.ExpectQuery("SELECT").
				WithArgs(tt.year, tt.month).
				WillReturnRows(tt.mockRows)

			got, err := repo.AggregateByEmployee(ctx, tt.year, tt.month)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(got))
				if tt.wantCount > 0 {
					total := 0
					for _, r := range got {
						total += r.TotalAmount
					}
					assert.Equal(t, tt.wantTotal, total)
				}
			}

			// 全ての期待値が満たされたか確認
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReportRepository_AggregateByCategory(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewReportRepository(mock)
	ctx := context.Background()

	mockRows := pgxmock.NewRows([]string{"category_id", "category_name", "count", "total_amount"}).
		AddRow("c1", "交通費", 5, 10000)

	mock.ExpectQuery("SELECT").
		WithArgs(2026, 4).
		WillReturnRows(mockRows)

	got, err := repo.AggregateByCategory(ctx, 2026, 4)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "交通費", got[0].CategoryName)
	assert.Equal(t, 10000, got[0].TotalAmount)

	assert.NoError(t, mock.ExpectationsWereMet())
}
