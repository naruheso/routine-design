package model

import "time"

// Expense は経費エンティティを表す（集約ルート）。
type Expense struct {
	ID               string     `json:"id"`
	UserID           string     `json:"userId"`
	UserName         string     `json:"userName,omitempty"`
	ExpenseDate      string     `json:"expenseDate"` // YYYY-MM-DD
	CategoryID       string     `json:"categoryId"`
	CategoryName     string     `json:"categoryName,omitempty"`
	Amount           int        `json:"amount"`
	Description      string     `json:"description"`
	ReceiptImagePath *string    `json:"receiptImagePath"`
	Status           string     `json:"status"`
	SubmittedAt      *time.Time `json:"submittedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// IsEditable は経費が編集可能か判定する（BR-01）。
func (e *Expense) IsEditable() bool {
	return e.Status == StatusDraft || e.Status == StatusReturned
}

// CreateExpenseRequest は経費作成リクエストを表す。
type CreateExpenseRequest struct {
	ExpenseDate string `json:"expenseDate" binding:"required"`
	CategoryID  string `json:"categoryId" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=1,max=999999"`
	Description string `json:"description" binding:"required,max=500"`
}

// UpdateExpenseRequest は経費更新リクエストを表す。
type UpdateExpenseRequest struct {
	ExpenseDate string `json:"expenseDate" binding:"required"`
	CategoryID  string `json:"categoryId" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=1,max=999999"`
	Description string `json:"description" binding:"required,max=500"`
}

// ステータス定数
const (
	StatusDraft                  = "draft"
	StatusPendingManager         = "pending_manager"
	StatusPendingExpenseAdmin    = "pending_expense_admin"
	StatusPendingFinanceDirector = "pending_finance_director"
	StatusReturned               = "returned"
	StatusRejected               = "rejected"
	StatusApproved               = "approved"
)

// NextApprovalStatus は現在のステータスに対する承認後の次ステータスを返す。
func NextApprovalStatus(current string) string {
	switch current {
	case StatusPendingManager:
		return StatusPendingExpenseAdmin
	case StatusPendingExpenseAdmin:
		return StatusPendingFinanceDirector
	case StatusPendingFinanceDirector:
		return StatusApproved
	default:
		return ""
	}
}

// RequiredRoleForStatus は承認待ちステータスに対応する承認者ロールを返す。
func RequiredRoleForStatus(status string) string {
	switch status {
	case StatusPendingManager:
		return RoleManager
	case StatusPendingExpenseAdmin:
		return RoleExpenseAdmin
	case StatusPendingFinanceDirector:
		return RoleFinanceDirector
	default:
		return ""
	}
}

// StatusDisplayName はステータスの日本語表示名を返す。
func StatusDisplayName(status string) string {
	names := map[string]string{
		StatusDraft:                  "下書き",
		StatusPendingManager:         "上長承認待ち",
		StatusPendingExpenseAdmin:    "経費担当承認待ち",
		StatusPendingFinanceDirector: "経理部長承認待ち",
		StatusReturned:               "差し戻し",
		StatusRejected:               "否認",
		StatusApproved:               "承認完了",
	}
	if name, ok := names[status]; ok {
		return name
	}
	return status
}
