package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/amical/routine-design/backend/internal/model"
)

// ExpenseRepository は経費データのDB操作を行う。
type ExpenseRepository struct {
	pool *pgxpool.Pool
}

// NewExpenseRepository はExpenseRepositoryを生成する。
func NewExpenseRepository(pool *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{pool: pool}
}

// FindByID はIDで経費を検索する。
func (r *ExpenseRepository) FindByID(ctx context.Context, id string) (*model.Expense, error) {
	query := `
		SELECT e.id, e.user_id, u.name, e.expense_date, e.category_id, c.name,
		       e.amount, e.description, e.receipt_image_path, e.status,
		       e.submitted_at, e.created_at, e.updated_at
		FROM expenses e
		JOIN users u ON e.user_id = u.id
		JOIN categories c ON e.category_id = c.id
		WHERE e.id = $1`

	return r.scanExpense(r.pool.QueryRow(ctx, query, id))
}

// FindByUserID はユーザーIDで経費一覧を取得する。
func (r *ExpenseRepository) FindByUserID(ctx context.Context, userID string) ([]model.Expense, error) {
	query := `
		SELECT e.id, e.user_id, u.name, e.expense_date, e.category_id, c.name,
		       e.amount, e.description, e.receipt_image_path, e.status,
		       e.submitted_at, e.created_at, e.updated_at
		FROM expenses e
		JOIN users u ON e.user_id = u.id
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = $1
		ORDER BY e.created_at DESC`

	return r.queryExpenses(ctx, query, userID)
}

// FindByStatuses はステータス群に該当する経費一覧を取得する。
func (r *ExpenseRepository) FindByStatuses(ctx context.Context, statuses []string) ([]model.Expense, error) {
	query := `
		SELECT e.id, e.user_id, u.name, e.expense_date, e.category_id, c.name,
		       e.amount, e.description, e.receipt_image_path, e.status,
		       e.submitted_at, e.created_at, e.updated_at
		FROM expenses e
		JOIN users u ON e.user_id = u.id
		JOIN categories c ON e.category_id = c.id
		WHERE e.status = ANY($1)
		ORDER BY e.submitted_at DESC NULLS LAST`

	return r.queryExpenses(ctx, query, statuses)
}

// Create は経費を新規作成する。
func (r *ExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	query := `
		INSERT INTO expenses (user_id, expense_date, category_id, amount, description, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		expense.UserID, expense.ExpenseDate, expense.CategoryID,
		expense.Amount, expense.Description, expense.Status,
	).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
}

// Update は経費を更新する。
func (r *ExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	query := `
		UPDATE expenses
		SET expense_date = $2, category_id = $3, amount = $4, description = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status IN ('draft', 'returned')
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		expense.ID, expense.ExpenseDate, expense.CategoryID,
		expense.Amount, expense.Description,
	).Scan(&expense.UpdatedAt)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("expense %s: %w", expense.ID, fmt.Errorf("この経費データは現在編集できないステータスです"))
	}
	return err
}

// UpdateStatus はステータスを更新する（排他制御付き）。
func (r *ExpenseRepository) UpdateStatus(ctx context.Context, id string, expectedStatus string, newStatus string) error {
	now := time.Now()
	query := `
		UPDATE expenses
		SET status = $3, submitted_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = $2
		RETURNING id`

	var submittedAt *time.Time
	if newStatus == model.StatusPendingManager {
		submittedAt = &now
	}

	var returnedID string
	err := r.pool.QueryRow(ctx, query, id, expectedStatus, newStatus, submittedAt).Scan(&returnedID)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("排他制御エラー: %w", fmt.Errorf("この申請は既に他のユーザーによって処理されています"))
	}
	return err
}

// Delete は下書きの経費を削除する。
func (r *ExpenseRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2 AND status = 'draft'`
	tag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("expense not found or not deletable")
	}
	return nil
}

func (r *ExpenseRepository) scanExpense(row pgx.Row) (*model.Expense, error) {
	var e model.Expense
	var expenseDate time.Time
	err := row.Scan(
		&e.ID, &e.UserID, &e.UserName, &expenseDate, &e.CategoryID, &e.CategoryName,
		&e.Amount, &e.Description, &e.ReceiptImagePath, &e.Status,
		&e.SubmittedAt, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan expense: %w", err)
	}
	e.ExpenseDate = expenseDate.Format("2006-01-02")
	return &e, nil
}

func (r *ExpenseRepository) queryExpenses(ctx context.Context, query string, args ...interface{}) ([]model.Expense, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query expenses: %w", err)
	}
	defer rows.Close()

	var expenses = make([]model.Expense, 0)
	for rows.Next() {
		var e model.Expense
		var expenseDate time.Time
		err := rows.Scan(
			&e.ID, &e.UserID, &e.UserName, &expenseDate, &e.CategoryID, &e.CategoryName,
			&e.Amount, &e.Description, &e.ReceiptImagePath, &e.Status,
			&e.SubmittedAt, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan expense row: %w", err)
		}
		e.ExpenseDate = expenseDate.Format("2006-01-02")
		expenses = append(expenses, e)
	}
	return expenses, nil
}
