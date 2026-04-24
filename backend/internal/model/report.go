package model

// EmployeeReport は社員別月別集計レポートの1行を表す。
type EmployeeReport struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	Count       int    `json:"count"`
	TotalAmount int    `json:"totalAmount"`
}

// CategoryReport は勘定科目別集計レポートの1行を表す。
type CategoryReport struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Count        int     `json:"count"`
	TotalAmount  int     `json:"totalAmount"`
	Percentage   float64 `json:"percentage"`
}

// ReportSummary は集計レポート全体のサマリーを表す。
type ReportSummary struct {
	Year       int `json:"year"`
	Month      int `json:"month"`
	TotalCount int `json:"totalCount"`
	GrandTotal int `json:"grandTotal"`
}

// EmployeeReportResponse は社員別集計レポートのAPIレスポンス。
type EmployeeReportResponse struct {
	Summary   ReportSummary    `json:"summary"`
	Employees []EmployeeReport `json:"employees"`
}

// CategoryReportResponse は勘定科目別集計レポートのAPIレスポンス。
type CategoryReportResponse struct {
	Summary    ReportSummary    `json:"summary"`
	Categories []CategoryReport `json:"categories"`
}
