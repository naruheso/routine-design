package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/service"
)

// ReportHandler は集計レポートAPIハンドラを提供する。
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler はReportHandlerを生成する。
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// ByEmployee は社員別月別集計レポートを取得する。
// GET /api/reports/by-employee?year=2026&month=4
func (h *ReportHandler) ByEmployee(c *gin.Context) {
	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.reportService.GetEmployeeReport(c.Request.Context(), year, month)
	if err != nil {
		handleReportError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Data: result})
}

// ByCategory は勘定科目別集計レポートを取得する。
// GET /api/reports/by-category?year=2026&month=4
func (h *ReportHandler) ByCategory(c *gin.Context) {
	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.reportService.GetCategoryReport(c.Request.Context(), year, month)
	if err != nil {
		handleReportError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Data: result})
}

// parseYearMonth はクエリパラメータからyear, monthを取得する。
// 未指定の場合は現在の年月をデフォルトとする。
func parseYearMonth(c *gin.Context) (int, int, error) {
	now := time.Now()

	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, fmt.Errorf("year パラメータが不正です")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, 0, fmt.Errorf("month パラメータが不正です")
	}

	return year, month, nil
}

// handleReportError はレポートサービスのエラーをHTTPレスポンスに変換する。
func handleReportError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrValidation) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "内部エラー"})
}
