package service

import (
	"context"
	"fmt"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
)

// ExpenseService は経費のビジネスロジックを提供する。
type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
	approvalRepo *repository.ApprovalRepository
}

// NewExpenseService はExpenseServiceを生成する。
func NewExpenseService(expenseRepo *repository.ExpenseRepository, approvalRepo *repository.ApprovalRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo, approvalRepo: approvalRepo}
}

// FindByID は経費を取得する。
func (s *ExpenseService) FindByID(ctx context.Context, id string) (*model.Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find expense %s: %w", id, err)
	}
	if expense == nil {
		return nil, fmt.Errorf("expense %s: %w", id, ErrNotFound)
	}
	return expense, nil
}

// FindByUserID は指定ユーザーの経費一覧を取得する。
func (s *ExpenseService) FindByUserID(ctx context.Context, userID string) ([]model.Expense, error) {
	expenses, err := s.expenseRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find expenses by user %s: %w", userID, err)
	}
	return expenses, nil
}

// Create は経費を新規作成する（下書きステータス）。
func (s *ExpenseService) Create(ctx context.Context, userID string, req model.CreateExpenseRequest) (*model.Expense, error) {
	expense := &model.Expense{
		UserID:      userID,
		ExpenseDate: req.ExpenseDate,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Status:      model.StatusDraft,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("create expense: %w", err)
	}

	return s.expenseRepo.FindByID(ctx, expense.ID)
}

// Update は経費を更新する（下書きまたは差し戻しのみ）。
func (s *ExpenseService) Update(ctx context.Context, id string, userID string, req model.UpdateExpenseRequest) (*model.Expense, error) {
	expense, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 権限チェック: 自分の経費のみ編集可能
	if expense.UserID != userID {
		return nil, fmt.Errorf("expense %s: %w", id, ErrForbidden)
	}

	// ビジネスルール BR-01: 修正可能条件
	if !expense.IsEditable() {
		return nil, fmt.Errorf("この経費データは現在編集できないステータスです: %w", ErrValidation)
	}

	expense.ExpenseDate = req.ExpenseDate
	expense.CategoryID = req.CategoryID
	expense.Amount = req.Amount
	expense.Description = req.Description

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("update expense: %w", err)
	}

	return s.expenseRepo.FindByID(ctx, id)
}

// Submit は経費を申請する（draft/returned → pending_manager）。
func (s *ExpenseService) Submit(ctx context.Context, id string, userID string) (*model.Expense, error) {
	expense, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 権限チェック: 自分の経費のみ申請可能
	if expense.UserID != userID {
		return nil, fmt.Errorf("expense %s: %w", id, ErrForbidden)
	}

	// 状態遷移チェック: draft or returned のみ
	if expense.Status != model.StatusDraft && expense.Status != model.StatusReturned {
		return nil, fmt.Errorf("この経費は申請できないステータスです: %w", ErrValidation)
	}

	// ステータス更新
	if err := s.expenseRepo.UpdateStatus(ctx, id, expense.Status, model.StatusPendingManager); err != nil {
		return nil, fmt.Errorf("submit expense: %w", err)
	}

	// 承認履歴に「申請」を記録
	action := model.ActionSubmit
	history := &model.ApprovalHistory{
		ExpenseID: id,
		Action:    action,
		ActorID:   userID,
		ActorRole: model.RoleApplicant,
	}
	if err := s.approvalRepo.Create(ctx, history); err != nil {
		return nil, fmt.Errorf("create submit history: %w", err)
	}

	return s.expenseRepo.FindByID(ctx, id)
}

// Delete は下書きの経費を削除する。
func (s *ExpenseService) Delete(ctx context.Context, id string, userID string) error {
	return s.expenseRepo.Delete(ctx, id, userID)
}

// GetHistories は経費の承認履歴を取得する。
func (s *ExpenseService) GetHistories(ctx context.Context, expenseID string) ([]model.ApprovalHistory, error) {
	return s.approvalRepo.FindByExpenseID(ctx, expenseID)
}
