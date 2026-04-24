package service

import (
	"context"
	"fmt"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
)

// ApprovalService は承認ワークフローのビジネスロジックを提供する。
type ApprovalService struct {
	expenseRepo  *repository.ExpenseRepository
	approvalRepo *repository.ApprovalRepository
}

// NewApprovalService はApprovalServiceを生成する。
func NewApprovalService(expenseRepo *repository.ExpenseRepository, approvalRepo *repository.ApprovalRepository) *ApprovalService {
	return &ApprovalService{expenseRepo: expenseRepo, approvalRepo: approvalRepo}
}

// GetPendingForRole は指定ロールの承認待ち経費一覧を取得する。
func (s *ApprovalService) GetPendingForRole(ctx context.Context, roles []string) ([]model.Expense, error) {
	var statuses []string
	for _, role := range roles {
		switch role {
		case model.RoleManager:
			statuses = append(statuses, model.StatusPendingManager)
		case model.RoleExpenseAdmin:
			statuses = append(statuses, model.StatusPendingExpenseAdmin)
		case model.RoleFinanceDirector:
			statuses = append(statuses, model.StatusPendingFinanceDirector)
		}
	}
	if len(statuses) == 0 {
		return []model.Expense{}, nil
	}
	return s.expenseRepo.FindByStatuses(ctx, statuses)
}

// Approve は経費を承認する。
func (s *ApprovalService) Approve(ctx context.Context, expenseID string, actorID string, actorRoles []string, comment *string) error {
	return s.processApproval(ctx, expenseID, actorID, actorRoles, model.ActionApprove, comment)
}

// Return は経費を差し戻す。
func (s *ApprovalService) Return(ctx context.Context, expenseID string, actorID string, actorRoles []string, comment *string) error {
	// BR-02: 差し戻し時コメント必須
	if comment == nil || *comment == "" {
		return fmt.Errorf("差し戻しの理由をコメントに入力してください: %w", ErrValidation)
	}
	return s.processApproval(ctx, expenseID, actorID, actorRoles, model.ActionReturn, comment)
}

// Reject は経費を否認する。
func (s *ApprovalService) Reject(ctx context.Context, expenseID string, actorID string, actorRoles []string, comment *string) error {
	// BR-02: 否認時コメント必須
	if comment == nil || *comment == "" {
		return fmt.Errorf("否認の理由をコメントに入力してください: %w", ErrValidation)
	}
	return s.processApproval(ctx, expenseID, actorID, actorRoles, model.ActionReject, comment)
}

func (s *ApprovalService) processApproval(ctx context.Context, expenseID string, actorID string, actorRoles []string, action string, comment *string) error {
	// 経費取得
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("find expense: %w", err)
	}
	if expense == nil {
		return fmt.Errorf("expense %s: %w", expenseID, ErrNotFound)
	}

	// BR-07: 自己承認の禁止
	if expense.UserID == actorID {
		return fmt.Errorf("自分の申請は承認できません: %w", ErrForbidden)
	}

	// BR-03: 承認者のロールと現在のステータスが一致するか検証
	requiredRole := model.RequiredRoleForStatus(expense.Status)
	if requiredRole == "" {
		return fmt.Errorf("この申請は承認できないステータスです: %w", ErrValidation)
	}

	hasRequiredRole := false
	var matchedRole string
	for _, r := range actorRoles {
		if r == requiredRole {
			hasRequiredRole = true
			matchedRole = r
			break
		}
	}
	if !hasRequiredRole {
		return fmt.Errorf("この承認段階に対する権限がありません: %w", ErrForbidden)
	}

	// 次のステータスを決定
	var newStatus string
	switch action {
	case model.ActionApprove:
		newStatus = model.NextApprovalStatus(expense.Status)
	case model.ActionReturn:
		newStatus = model.StatusReturned
	case model.ActionReject:
		newStatus = model.StatusRejected
	}

	// ステータス更新（排他制御付き）
	if err := s.expenseRepo.UpdateStatus(ctx, expenseID, expense.Status, newStatus); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	// 承認履歴を記録
	history := &model.ApprovalHistory{
		ExpenseID: expenseID,
		Action:    action,
		ActorID:   actorID,
		ActorRole: matchedRole,
		Comment:   comment,
	}
	if err := s.approvalRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("create approval history: %w", err)
	}

	return nil
}
