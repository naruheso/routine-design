package model

// SuccessResponse は成功レスポンスを生成する。
func SuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}

// SuccessResponseWithCount は件数付き成功レスポンスを生成する。
func SuccessResponseWithCount(data interface{}, count int) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
		"count":   count,
	}
}

// SuccessResponseWithPagination はページネーション付き成功レスポンスを生成する。
func SuccessResponseWithPagination(data interface{}, pagination Pagination) map[string]interface{} {
	return map[string]interface{}{
		"success":    true,
		"data":       data,
		"pagination": pagination,
	}
}

// ErrorResponse はエラーレスポンスを生成する。
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   message,
	}
}

// Pagination はページネーション情報を表す。
type Pagination struct {
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	TotalCount  int  `json:"totalCount"`
	Limit       int  `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrev     bool `json:"hasPrev"`
}
