export interface EmployeeReport {
  userId: string
  userName: string
  count: number
  totalAmount: number
}

export interface CategoryReport {
  categoryId: string
  categoryName: string
  count: number
  totalAmount: number
  percentage: number
}

export interface ReportSummary {
  year: number
  month: number
  totalCount: number
  grandTotal: number
}

export interface EmployeeReportResponse {
  summary: ReportSummary
  employees: EmployeeReport[]
}

export interface CategoryReportResponse {
  summary: ReportSummary
  categories: CategoryReport[]
}
