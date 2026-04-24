import { ref } from 'vue'
import { defineStore } from 'pinia'
import { fetchPendingApprovals, approveExpense, returnExpense, rejectExpense } from '@/api/approvals'
import type { Expense, ApprovalRequest } from '@/types/expense'

export const useApprovalStore = defineStore('approval', () => {
  // State
  const pendingExpenses = ref<Expense[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function loadPending() {
    loading.value = true
    error.value = null
    try {
      const result = await fetchPendingApprovals()
      if (result.success && result.data !== undefined) {
        pendingExpenses.value = result.data
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  async function approve(id: string, comment?: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const req: ApprovalRequest = comment ? { comment } : {}
      const result = await approveExpense(id, req)
      if (result.success) {
        await loadPending()
        return true
      } else {
        error.value = result.error ?? '承認に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function returnBack(id: string, comment: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await returnExpense(id, { comment })
      if (result.success) {
        await loadPending()
        return true
      } else {
        error.value = result.error ?? '差し戻しに失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  async function reject(id: string, comment: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await rejectExpense(id, { comment })
      if (result.success) {
        await loadPending()
        return true
      } else {
        error.value = result.error ?? '否認に失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    pendingExpenses, loading, error,
    loadPending, approve, returnBack, reject,
  }
})
