package repository

import (
	"context"

	"github.com/amical/routine-design/backend/internal/model"
)

// ExpenseRepository は経費データのデータアクセスインターフェース。
type ExpenseRepository interface {
	FindByID(ctx context.Context, id string) (*model.Expense, error)
	FindByUserID(ctx context.Context, userID string) ([]model.Expense, error)
	FindByStatuses(ctx context.Context, statuses []string) ([]model.Expense, error)
	Create(ctx context.Context, expense *model.Expense) error
	Update(ctx context.Context, expense *model.Expense) error
	UpdateStatus(ctx context.Context, id string, expectedStatus string, newStatus string) error
	Delete(ctx context.Context, id string, userID string) error
}

// ApprovalRepository は承認履歴のデータアクセスインターフェース。
type ApprovalRepository interface {
	Create(ctx context.Context, history *model.ApprovalHistory) error
	FindByExpenseID(ctx context.Context, expenseID string) ([]model.ApprovalHistory, error)
}

// UserRepository はユーザーデータのデータアクセスインターフェース。
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	GetRoles(ctx context.Context, userID string) ([]string, error)
}

// CategoryRepository は勘定科目のデータアクセスインターフェース。
type CategoryRepository interface {
	FindAllActive(ctx context.Context) ([]model.Category, error)
}

// ReportRepository はレポート集計のデータアクセスインターフェース。
type ReportRepository interface {
	AggregateByEmployee(ctx context.Context, year, month int) ([]model.EmployeeReport, error)
	AggregateByCategory(ctx context.Context, year, month int) ([]model.CategoryReport, error)
}
