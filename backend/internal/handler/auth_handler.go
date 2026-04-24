package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/amical/routine-design/backend/internal/model"
	"github.com/amical/routine-design/backend/internal/service"
)

// handleError は共通エラーハンドラ。Service層のエラーをHTTPステータスに変換する。
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, model.ErrorResponse("not found"))
	case errors.Is(err, service.ErrValidation):
		c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, model.ErrorResponse(err.Error()))
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(err.Error()))
	case errors.Is(err, service.ErrConflict), errors.Is(err, service.ErrStatusConflict):
		c.JSON(http.StatusConflict, model.ErrorResponse(err.Error()))
	default:
		slog.Error("internal error", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("システムエラーが発生しました。時間をおいて再度お試しください。"))
	}
}

// AuthHandler は認証関連のHTTPハンドラを提供する。
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler はAuthHandlerを生成する。
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login はログインエンドポイント。
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("メールアドレスとパスワードを入力してください"))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}
