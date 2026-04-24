import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchPendingApprovals, approveExpense, returnExpense, rejectExpense } from '@/api/approvals'

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '@/api/client'

describe('approvals API', () => {
  beforeEach(() => { vi.clearAllMocks() })

  describe('fetchPendingApprovals()', () => {
    it('正常系: 承認待ち一覧を取得する', async () => {
      const mockExpenses = [
        { id: 'exp-1', userName: '田中太郎', amount: 1500, status: 'pending_manager' },
      ]
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: mockExpenses })
      const result = await fetchPendingApprovals()
      expect(apiClient.get).toHaveBeenCalledWith('/approvals/pending')
      expect(result.success).toBe(true)
      expect(result.data).toHaveLength(1)
    })

    it('正常系: 承認待ちが0件の場合', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: [] })
      const result = await fetchPendingApprovals()
      expect(result.data).toEqual([])
    })

    it('異常系: 権限不足で取得失敗', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: false, error: '権限がありません' })
      const result = await fetchPendingApprovals()
      expect(result.success).toBe(false)
    })
  })

  describe('approveExpense()', () => {
    it('正常系: コメント付きで承認する', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: true })
      const result = await approveExpense('exp-1', { comment: '問題ありません' })
      expect(apiClient.post).toHaveBeenCalledWith('/approvals/exp-1/approve', { comment: '問題ありません' })
      expect(result.success).toBe(true)
    })

    it('正常系: コメントなしで承認する', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: true })
      const result = await approveExpense('exp-1')
      expect(apiClient.post).toHaveBeenCalledWith('/approvals/exp-1/approve', undefined)
      expect(result.success).toBe(true)
    })
  })

  describe('returnExpense()', () => {
    it('正常系: コメント付きで差し戻す', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: true })
      const result = await returnExpense('exp-1', { comment: '領収書を添付してください' })
      expect(apiClient.post).toHaveBeenCalledWith('/approvals/exp-1/return', { comment: '領収書を添付してください' })
      expect(result.success).toBe(true)
    })

    it('異常系: 差し戻し失敗', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: false, error: '処理に失敗しました' })
      const result = await returnExpense('exp-1', { comment: '理由' })
      expect(result.success).toBe(false)
    })
  })

  describe('rejectExpense()', () => {
    it('正常系: コメント付きで否認する', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: true })
      const result = await rejectExpense('exp-1', { comment: '経費対象外です' })
      expect(apiClient.post).toHaveBeenCalledWith('/approvals/exp-1/reject', { comment: '経費対象外です' })
      expect(result.success).toBe(true)
    })
  })
})
