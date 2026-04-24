package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/amical/routine-design/backend/internal/model"
)

// pgUserRepository はユーザーデータのDB操作を行うPostgreSQL実装。
type pgUserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository はUserRepositoryインターフェースを満たす構造体を生成する。
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepository{pool: pool}
}

// FindByEmail はメールアドレスでユーザーを検索する。
func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	// ロール取得
	roles, err := r.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("find roles for user %s: %w", user.ID, err)
	}
	user.Roles = roles

	return &user, nil
}

// FindByID はIDでユーザーを検索する。
func (r *pgUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	roles, err := r.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("find roles for user %s: %w", user.ID, err)
	}
	user.Roles = roles

	return &user, nil
}

// GetRoles はユーザーのロール一覧を取得する。
func (r *pgUserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT role FROM user_roles WHERE user_id = $1`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
