package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/amical/routine-design/backend/internal/model"
)

// pgApprovalRepository は承認履歴のDB操作を行うPostgreSQL実装。
type pgApprovalRepository struct {
	pool *pgxpool.Pool
}

// NewApprovalRepository はApprovalRepositoryインターフェースを満たす構造体を生成する。
func NewApprovalRepository(pool *pgxpool.Pool) ApprovalRepository {
	return &pgApprovalRepository{pool: pool}
}

// Create は承認履歴を記録する。
func (r *pgApprovalRepository) Create(ctx context.Context, history *model.ApprovalHistory) error {
	query := `
		INSERT INTO approval_histories (expense_id, action, actor_id, actor_role, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return GetDBQuerier(ctx, r.pool).QueryRow(ctx, query,
		history.ExpenseID, history.Action, history.ActorID, history.ActorRole, history.Comment,
	).Scan(&history.ID, &history.CreatedAt)
}

// FindByExpenseID は経費IDで承認履歴を取得する。
func (r *pgApprovalRepository) FindByExpenseID(ctx context.Context, expenseID string) ([]model.ApprovalHistory, error) {
	query := `
		SELECT ah.id, ah.expense_id, ah.action, ah.actor_id, u.name, ah.actor_role, ah.comment, ah.created_at
		FROM approval_histories ah
		JOIN users u ON ah.actor_id = u.id
		WHERE ah.expense_id = $1
		ORDER BY ah.created_at ASC`

	rows, err := GetDBQuerier(ctx, r.pool).Query(ctx, query, expenseID)
	if err != nil {
		return nil, fmt.Errorf("query approval histories: %w", err)
	}
	defer rows.Close()

	var histories = make([]model.ApprovalHistory, 0)
	for rows.Next() {
		var h model.ApprovalHistory
		err := rows.Scan(&h.ID, &h.ExpenseID, &h.Action, &h.ActorID, &h.ActorName, &h.ActorRole, &h.Comment, &h.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan approval history: %w", err)
		}
		histories = append(histories, h)
	}
	return histories, nil
}

// pgCategoryRepository は勘定科目のDB操作を行うPostgreSQL実装。
type pgCategoryRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryRepository はCategoryRepositoryインターフェースを満たす構造体を生成する。
func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &pgCategoryRepository{pool: pool}
}

// FindAllActive は有効な勘定科目を取得する。
func (r *pgCategoryRepository) FindAllActive(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, code, name, sort_order, is_active FROM categories WHERE is_active = true ORDER BY sort_order`
	rows, err := GetDBQuerier(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	var categories = make([]model.Category, 0)
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.SortOrder, &c.IsActive); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}
