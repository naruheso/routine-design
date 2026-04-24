package service

import (
	"context"
	"fmt"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
)

// ReportService は集計レポートのビジネスロジックを提供する。
type ReportService struct {
	reportRepo repository.ReportRepository
}

// NewReportService はReportServiceを生成する。
func NewReportService(reportRepo repository.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

// GetEmployeeReport は社員別月別集計レポートを取得する。
func (s *ReportService) GetEmployeeReport(ctx context.Context, year, month int) (*model.EmployeeReportResponse, error) {
	if err := validateYearMonth(year, month); err != nil {
		return nil, err
	}

	employees, err := s.reportRepo.AggregateByEmployee(ctx, year, month)
	if err != nil {
		return nil, fmt.Errorf("aggregate by employee: %w", err)
	}

	// サマリー計算
	totalCount := 0
	grandTotal := 0
	for _, e := range employees {
		totalCount += e.Count
		grandTotal += e.TotalAmount
	}

	return &model.EmployeeReportResponse{
		Summary: model.ReportSummary{
			Year:       year,
			Month:      month,
			TotalCount: totalCount,
			GrandTotal: grandTotal,
		},
		Employees: employees,
	}, nil
}

// GetCategoryReport は勘定科目別集計レポートを取得する。
func (s *ReportService) GetCategoryReport(ctx context.Context, year, month int) (*model.CategoryReportResponse, error) {
	if err := validateYearMonth(year, month); err != nil {
		return nil, err
	}

	categories, err := s.reportRepo.AggregateByCategory(ctx, year, month)
	if err != nil {
		return nil, fmt.Errorf("aggregate by category: %w", err)
	}

	// 総合計と構成比の計算
	grandTotal := 0
	totalCount := 0
	for _, c := range categories {
		grandTotal += c.TotalAmount
		totalCount += c.Count
	}

	// 構成比をパーセンテージとして設定
	if grandTotal > 0 {
		for i := range categories {
			categories[i].Percentage = float64(categories[i].TotalAmount) / float64(grandTotal) * 100.0
		}
	}

	return &model.CategoryReportResponse{
		Summary: model.ReportSummary{
			Year:       year,
			Month:      month,
			TotalCount: totalCount,
			GrandTotal: grandTotal,
		},
		Categories: categories,
	}, nil
}

// validateYearMonth は年月のバリデーションを行う。
func validateYearMonth(year, month int) error {
	if year < 2020 || year > 2100 {
		return fmt.Errorf("年の指定が不正です（2020〜2100）: %w", ErrValidation)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("月の指定が不正です（1〜12）: %w", ErrValidation)
	}
	return nil
}
