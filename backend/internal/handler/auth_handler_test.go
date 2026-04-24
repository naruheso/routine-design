package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Login_BadRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "異常系: 空のリクエストボディ",
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantError:  "メールアドレスとパスワードを入力してください",
		},
		{
			name:       "異常系: 不正なJSON",
			body:       "{invalid json}",
			wantStatus: http.StatusBadRequest,
			wantError:  "メールアドレスとパスワードを入力してください",
		},
		{
			name:       "異常系: 空オブジェクト",
			body:       "{}",
			wantStatus: http.StatusBadRequest,
			wantError:  "メールアドレスとパスワードを入力してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// AuthServiceがnilなのでbindJSON段階でエラーにならないと先に進んでpanicするが、
			// 空ボディや不正JSONはbindJSONでエラーになるのでテスト可能
			h := &AuthHandler{authService: nil}
			h.Login(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var body map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			require.NoError(t, err)
			assert.Equal(t, false, body["success"])
			assert.Contains(t, body["error"], tt.wantError)
		})
	}
}

func TestAuthHandler_Login_ResponseFormat(t *testing.T) {
	t.Parallel()

	// レスポンス互換性チェック: successフィールドとerrorフィールドの存在
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h := &AuthHandler{authService: nil}
	h.Login(c)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	// レスポンス互換性: success フィールドが存在
	_, hasSuccess := body["success"]
	assert.True(t, hasSuccess, "レスポンスに 'success' フィールドが必要")

	// エラー時は error フィールドも存在
	_, hasError := body["error"]
	assert.True(t, hasError, "エラーレスポンスに 'error' フィールドが必要")
}
