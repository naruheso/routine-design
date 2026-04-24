package model

import "time"

// ApprovalHistory は承認履歴（不変エンティティ）を表す。
type ApprovalHistory struct {
	ID        string    `json:"id"`
	ExpenseID string    `json:"expenseId"`
	Action    string    `json:"action"`
	ActorID   string    `json:"actorId"`
	ActorName string    `json:"actorName,omitempty"`
	ActorRole string    `json:"actorRole"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

// ApprovalRequest は承認/差し戻し/否認リクエストを表す。
type ApprovalRequest struct {
	Comment *string `json:"comment"`
}

// アクション定数
const (
	ActionSubmit  = "submit"
	ActionApprove = "approve"
	ActionReturn  = "return"
	ActionReject  = "reject"
)

// Category は勘定科目エンティティを表す。
type Category struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
}
