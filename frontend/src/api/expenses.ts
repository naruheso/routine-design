import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { Expense, ExpenseDetail, CreateExpenseRequest, UpdateExpenseRequest, Category } from '@/types/expense'

export function fetchExpenses(): Promise<ApiResponse<Expense[]>> {
  return apiClient.get<Expense[]>('/expenses')
}

export function fetchExpenseById(id: string): Promise<ApiResponse<ExpenseDetail>> {
  return apiClient.get<ExpenseDetail>(`/expenses/${id}`)
}

export function createExpense(req: CreateExpenseRequest): Promise<ApiResponse<Expense>> {
  return apiClient.post<Expense>('/expenses', req)
}

export function updateExpense(id: string, req: UpdateExpenseRequest): Promise<ApiResponse<Expense>> {
  return apiClient.put<Expense>(`/expenses/${id}`, req)
}

export function deleteExpense(id: string): Promise<ApiResponse<void>> {
  return apiClient.delete<void>(`/expenses/${id}`)
}

export function submitExpense(id: string): Promise<ApiResponse<Expense>> {
  return apiClient.post<Expense>(`/expenses/${id}/submit`)
}

export function fetchCategories(): Promise<ApiResponse<Category[]>> {
  return apiClient.get<Category[]>('/categories')
}
