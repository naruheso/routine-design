package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/middleware"
	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/repository"
	"github.com/amical/routine-design/backend/internal/service"
)

// ApprovalHandler は承認関連のHTTPハンドラを提供する。
type ApprovalHandler struct {
	approvalService *service.ApprovalService
}

// NewApprovalHandler はApprovalHandlerを生成する。
func NewApprovalHandler(approvalService *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{approvalService: approvalService}
}

// GetPending は承認待ち経費一覧を取得する。
func (h *ApprovalHandler) GetPending(c *gin.Context) {
	roles := middleware.GetUserRoles(c)
	expenses, err := h.approvalService.GetPendingForRole(c.Request.Context(), roles)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(expenses))
}

// Approve は経費を承認する。
func (h *ApprovalHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	var req model.ApprovalRequest
	_ = c.ShouldBindJSON(&req)

	userID := middleware.GetUserID(c)
	roles := middleware.GetUserRoles(c)

	if err := h.approvalService.Approve(c.Request.Context(), id, userID, roles, req.Comment); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "承認しました"}))
}

// Return は経費を差し戻す。
func (h *ApprovalHandler) Return(c *gin.Context) {
	id := c.Param("id")
	var req model.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("リクエスト形式が不正です"))
		return
	}

	userID := middleware.GetUserID(c)
	roles := middleware.GetUserRoles(c)

	if err := h.approvalService.Return(c.Request.Context(), id, userID, roles, req.Comment); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "差し戻しました"}))
}

// Reject は経費を否認する。
func (h *ApprovalHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	var req model.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("リクエスト形式が不正です"))
		return
	}

	userID := middleware.GetUserID(c)
	roles := middleware.GetUserRoles(c)

	if err := h.approvalService.Reject(c.Request.Context(), id, userID, roles, req.Comment); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "否認しました"}))
}

// CategoryHandler は勘定科目関連のHTTPハンドラを提供する。
type CategoryHandler struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryHandler はCategoryHandlerを生成する。
func NewCategoryHandler(categoryRepo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

// List は有効な勘定科目一覧を取得する。
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryRepo.FindAllActive(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(categories))
}
