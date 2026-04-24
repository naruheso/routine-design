import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { Expense, ApprovalRequest } from '@/types/expense'

export function fetchPendingApprovals(): Promise<ApiResponse<Expense[]>> {
  return apiClient.get<Expense[]>('/approvals/pending')
}

export function approveExpense(id: string, req?: ApprovalRequest): Promise<ApiResponse<void>> {
  return apiClient.post<void>(`/approvals/${id}/approve`, req)
}

export function returnExpense(id: string, req: ApprovalRequest): Promise<ApiResponse<void>> {
  return apiClient.post<void>(`/approvals/${id}/return`, req)
}

export function rejectExpense(id: string, req: ApprovalRequest): Promise<ApiResponse<void>> {
  return apiClient.post<void>(`/approvals/${id}/reject`, req)
}
