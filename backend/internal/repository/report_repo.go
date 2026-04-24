package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/amical/routine-design/backend/internal/model"
)

// ReportRepository は集計レポートのDB操作を提供する。
type ReportRepository struct {
	pool *pgxpool.Pool
}

// NewReportRepository はReportRepositoryを生成する。
func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

// AggregateByEmployee は指定年月の承認済み経費を社員別に集計する（BR-05）。
func (r *ReportRepository) AggregateByEmployee(ctx context.Context, year, month int) ([]model.EmployeeReport, error) {
	query := `
		SELECT
			e.user_id,
			u.name AS user_name,
			COUNT(*) AS count,
			COALESCE(SUM(e.amount), 0) AS total_amount
		FROM expenses e
		JOIN users u ON u.id = e.user_id
		WHERE e.status = 'approved'
		  AND EXTRACT(YEAR FROM e.expense_date) = $1
		  AND EXTRACT(MONTH FROM e.expense_date) = $2
		GROUP BY e.user_id, u.name
		ORDER BY total_amount DESC
	`

	rows, err := r.pool.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("aggregate by employee: %w", err)
	}
	defer rows.Close()

	var reports = make([]model.EmployeeReport, 0)
	for rows.Next() {
		var rpt model.EmployeeReport
		if err := rows.Scan(&rpt.UserID, &rpt.UserName, &rpt.Count, &rpt.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan employee report: %w", err)
		}
		reports = append(reports, rpt)
	}
	return reports, nil
}

// AggregateByCategory は指定年月の承認済み経費を勘定科目別に集計する（BR-05）。
func (r *ReportRepository) AggregateByCategory(ctx context.Context, year, month int) ([]model.CategoryReport, error) {
	query := `
		SELECT
			e.category_id,
			c.name AS category_name,
			COUNT(*) AS count,
			COALESCE(SUM(e.amount), 0) AS total_amount
		FROM expenses e
		JOIN categories c ON c.id = e.category_id
		WHERE e.status = 'approved'
		  AND EXTRACT(YEAR FROM e.expense_date) = $1
		  AND EXTRACT(MONTH FROM e.expense_date) = $2
		GROUP BY e.category_id, c.name
		ORDER BY total_amount DESC
	`

	rows, err := r.pool.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("aggregate by category: %w", err)
	}
	defer rows.Close()

	var reports = make([]model.CategoryReport, 0)
	for rows.Next() {
		var rpt model.CategoryReport
		if err := rows.Scan(&rpt.CategoryID, &rpt.CategoryName, &rpt.Count, &rpt.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan category report: %w", err)
		}
		reports = append(reports, rpt)
	}
	return reports, nil
}
