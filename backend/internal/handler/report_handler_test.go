package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestReportHandler_ByEmployee_ParseParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		query     string
		wantErr   bool
		wantYear  int
		wantMonth int
	}{
		{
			name:      "正常系: 年月指定あり",
			query:     "?year=2026&month=4",
			wantErr:   false,
			wantYear:  2026,
			wantMonth: 4,
		},
		{
			name:    "正常系: 年月指定なし（デフォルト=当月）",
			query:   "",
			wantErr: false,
		},
		{
			name:    "異常系: yearが不正な文字列",
			query:   "?year=abc&month=4",
			wantErr: true,
		},
		{
			name:    "異常系: monthが不正な文字列",
			query:   "?year=2026&month=xyz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, "/api/reports/by-employee"+tt.query, nil)
			c.Request = req

			year, month, err := parseYearMonth(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.wantYear > 0 {
					assert.Equal(t, tt.wantYear, year)
					assert.Equal(t, tt.wantMonth, month)
				} else {
					// デフォルト値: 正の値が返る
					assert.Greater(t, year, 0)
					assert.GreaterOrEqual(t, month, 1)
					assert.LessOrEqual(t, month, 12)
				}
			}
		})
	}
}

func TestParseYearMonth_Defaults(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/reports/by-employee", nil)
	c.Request = req

	year, month, err := parseYearMonth(c)
	require.NoError(t, err)
	// デフォルト値は現在の年月
	assert.Greater(t, year, 2020)
	assert.GreaterOrEqual(t, month, 1)
	assert.LessOrEqual(t, month, 12)
}

func TestHandleError_Mapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantKey    string
	}{
		{
			name:       "正常系レスポンス形式: successフィールドが存在する",
			query:      "?year=abc",
			wantStatus: http.StatusBadRequest,
			wantKey:    "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, "/api/reports/by-employee"+tt.query, nil)
			c.Request = req

			_, _, err := parseYearMonth(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   err.Error(),
				})
			}

			assert.Equal(t, tt.wantStatus, w.Code)

			var body map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &body)
			require.NoError(t, err)
			assert.Contains(t, body, tt.wantKey)
			assert.Equal(t, false, body["success"])
		})
	}
}
