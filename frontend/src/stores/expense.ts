import { ref } from 'vue'
import { defineStore } from 'pinia'
import {
  fetchExpenses, fetchExpenseById, createExpense,
  updateExpense, deleteExpense, submitExpense, fetchCategories,
} from '@/api/expenses'
import type { Expense, ExpenseDetail, CreateExpenseRequest, UpdateExpenseRequest, Category } from '@/types/expense'

export const useExpenseStore = defineStore('expense', () => {
  // State
  const expenses = ref<Expense[]>([])
  const currentExpense = ref<ExpenseDetail | null>(null)
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function loadExpenses() {
    loading.value = true
    error.value = null
    try {
      const result = await fetchExpenses()
      if (result.success && result.data !== undefined) {
        expenses.value = result.data
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  async function loadExpenseDetail(id: string) {
    loading.value = true
    error.value = null
    try {
      const result = await fetchExpenseById(id)
      if (result.success && result.data) {
        currentExpense.value = result.data
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  async function addExpense(req: CreateExpenseRequest): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await createExpense(req)
      if (result.success) {
        await loadExpenses()
        return true
      } else {
        error.value = result.error ?? '作成に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function editExpense(id: string, req: UpdateExpenseRequest): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await updateExpense(id, req)
      if (result.success) {
        await loadExpenses()
        return true
      } else {
        error.value = result.error ?? '更新に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function removeExpense(id: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await deleteExpense(id)
      if (result.success) {
        await loadExpenses()
        return true
      } else {
        error.value = result.error ?? '削除に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function submit(id: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await submitExpense(id)
      if (result.success) {
        await loadExpenses()
        return true
      } else {
        error.value = result.error ?? '申請に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function loadCategories() {
    try {
      const result = await fetchCategories()
      if (result.success && result.data) {
        categories.value = result.data
      }
    } catch {
      // カテゴリ取得失敗は致命的でないため、エラーは表示しない
    }
  }

  return {
    expenses, currentExpense, categories, loading, error,
    loadExpenses, loadExpenseDetail, addExpense, editExpense,
    removeExpense, submit, loadCategories,
  }
})
