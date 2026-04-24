package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/middleware"
	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/service"
)

// ExpenseHandler は経費関連のHTTPハンドラを提供する。
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

// NewExpenseHandler はExpenseHandlerを生成する。
func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

// List は自分の経費一覧を取得する。
func (h *ExpenseHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	expenses, err := h.expenseService.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(expenses))
}

// GetByID は経費詳細を取得する。
func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	expense, err := h.expenseService.FindByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	// 承認履歴も取得
	histories, err := h.expenseService.GetHistories(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{
		"expense":   expense,
		"histories": histories,
	}))
}

// Create は経費を新規作成する。
func (h *ExpenseHandler) Create(c *gin.Context) {
	var req model.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
		return
	}

	userID := middleware.GetUserID(c)
	expense, err := h.expenseService.Create(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(expense))
}

// Update は経費を更新する。
func (h *ExpenseHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
		return
	}

	userID := middleware.GetUserID(c)
	expense, err := h.expenseService.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(expense))
}

// Submit は経費を申請する。
func (h *ExpenseHandler) Submit(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	expense, err := h.expenseService.Submit(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(expense))
}

// Delete は経費を削除する（下書きのみ）。
func (h *ExpenseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.expenseService.Delete(c.Request.Context(), id, userID); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
