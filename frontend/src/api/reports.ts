import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { EmployeeReportResponse, CategoryReportResponse } from '@/types/report'

export function fetchEmployeeReport(year: number, month: number): Promise<ApiResponse<EmployeeReportResponse>> {
  return apiClient.get<EmployeeReportResponse>(`/reports/by-employee?year=${year}&month=${month}`)
}

export function fetchCategoryReport(year: number, month: number): Promise<ApiResponse<CategoryReportResponse>> {
  return apiClient.get<CategoryReportResponse>(`/reports/by-category?year=${year}&month=${month}`)
}
