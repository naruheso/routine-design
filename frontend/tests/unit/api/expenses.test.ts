import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  fetchExpenses, fetchExpenseById, createExpense,
  updateExpense, deleteExpense, submitExpense, fetchCategories,
} from '@/api/expenses'

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '@/api/client'

const mockExpense = {
  id: 'exp-1', userId: 'user-1', userName: '田中太郎',
  expenseDate: '2026-04-01', categoryId: 'cat-1', categoryName: '交通費',
  amount: 1500, description: '電車代', receiptImagePath: null,
  status: 'draft', submittedAt: null,
  createdAt: '2026-04-01T00:00:00Z', updatedAt: '2026-04-01T00:00:00Z',
}

describe('expenses API', () => {
  beforeEach(() => { vi.clearAllMocks() })

  describe('fetchExpenses()', () => {
    it('正常系: 経費一覧を取得する', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: [mockExpense] })
      const result = await fetchExpenses()
      expect(apiClient.get).toHaveBeenCalledWith('/expenses')
      expect(result.success).toBe(true)
      expect(result.data).toHaveLength(1)
    })

    it('正常系: 経費が0件の場合', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: [] })
      const result = await fetchExpenses()
      expect(result.data).toEqual([])
    })
  })

  describe('fetchExpenseById()', () => {
    it('正常系: 経費詳細を取得する', async () => {
      const mockDetail = { expense: mockExpense, histories: [] }
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: mockDetail })
      const result = await fetchExpenseById('exp-1')
      expect(apiClient.get).toHaveBeenCalledWith('/expenses/exp-1')
      expect(result.data?.expense.id).toBe('exp-1')
    })

    it('異常系: 存在しないIDで取得失敗', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: false, error: '経費データが見つかりません' })
      const result = await fetchExpenseById('non-existent')
      expect(result.success).toBe(false)
    })
  })

  describe('createExpense()', () => {
    it('正常系: 経費を新規作成する', async () => {
      const req = { expenseDate: '2026-04-10', categoryId: 'cat-1', amount: 3000, description: 'タクシー代' }
      vi.mocked(apiClient.post).mockResolvedValue({ success: true, data: { ...mockExpense, ...req } })
      const result = await createExpense(req)
      expect(apiClient.post).toHaveBeenCalledWith('/expenses', req)
      expect(result.success).toBe(true)
    })

    it('異常系: バリデーションエラー', async () => {
      const req = { expenseDate: '', categoryId: '', amount: -100, description: '' }
      vi.mocked(apiClient.post).mockResolvedValue({ success: false, error: '金額は1円以上で入力してください' })
      const result = await createExpense(req)
      expect(result.success).toBe(false)
    })
  })

  describe('updateExpense()', () => {
    it('正常系: 経費を更新する', async () => {
      const req = { expenseDate: '2026-04-10', categoryId: 'cat-2', amount: 5000, description: '更新' }
      vi.mocked(apiClient.put).mockResolvedValue({ success: true, data: { ...mockExpense, ...req } })
      const result = await updateExpense('exp-1', req)
      expect(apiClient.put).toHaveBeenCalledWith('/expenses/exp-1', req)
      expect(result.success).toBe(true)
    })
  })

  describe('deleteExpense()', () => {
    it('正常系: 経費を削除する', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ success: true })
      const result = await deleteExpense('exp-1')
      expect(apiClient.delete).toHaveBeenCalledWith('/expenses/exp-1')
      expect(result.success).toBe(true)
    })

    it('異常系: 申請済み経費は削除不可', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({ success: false, error: '申請済みの経費は削除できません' })
      const result = await deleteExpense('exp-submitted')
      expect(result.success).toBe(false)
    })
  })

  describe('submitExpense()', () => {
    it('正常系: 経費を申請する', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ success: true, data: { ...mockExpense, status: 'pending_manager' } })
      const result = await submitExpense('exp-1')
      expect(apiClient.post).toHaveBeenCalledWith('/expenses/exp-1/submit')
      expect(result.success).toBe(true)
    })
  })

  describe('fetchCategories()', () => {
    it('正常系: カテゴリ一覧を取得する', async () => {
      const cats = [
        { id: 'cat-1', code: 'TRANSPORT', name: '交通費', sortOrder: 1, isActive: true },
        { id: 'cat-2', code: 'ENTERTAINMENT', name: '交際費', sortOrder: 2, isActive: true },
      ]
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: cats })
      const result = await fetchCategories()
      expect(apiClient.get).toHaveBeenCalledWith('/categories')
      expect(result.data).toHaveLength(2)
    })
  })
})
