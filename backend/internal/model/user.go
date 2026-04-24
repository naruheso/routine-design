package model

import "time"

// User はユーザーエンティティを表す。
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // JSONには出力しない
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// HasRole は指定ロールを保持しているか判定する。
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// LoginRequest はログインリクエストを表す。
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse はログインレスポンスを表す。
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ロール定数
const (
	RoleApplicant       = "applicant"
	RoleManager         = "manager"
	RoleExpenseAdmin    = "expense_admin"
	RoleFinanceDirector = "finance_director"
)
