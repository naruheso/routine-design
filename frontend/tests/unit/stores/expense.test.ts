import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useExpenseStore } from '@/stores/expense'

// APIモック
vi.mock('@/api/expenses', () => ({
  fetchExpenses: vi.fn(),
  fetchExpenseById: vi.fn(),
  createExpense: vi.fn(),
  updateExpense: vi.fn(),
  deleteExpense: vi.fn(),
  submitExpense: vi.fn(),
  fetchCategories: vi.fn(),
}))

import { fetchExpenses, createExpense, deleteExpense, submitExpense, fetchCategories } from '@/api/expenses'

describe('useExpenseStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('loadExpenses', () => {
    it('正常系: 経費一覧を取得', async () => {
      const mockExpenses = [
        { id: '1', userId: 'u1', expenseDate: '2026-04-20', categoryId: 'c1', categoryName: '交通費', amount: 1200, description: 'テスト', receiptImagePath: null, status: 'draft', submittedAt: null, createdAt: '', updatedAt: '' },
      ]
      vi.mocked(fetchExpenses).mockResolvedValue({
        success: true,
        data: mockExpenses,
      })

      const store = useExpenseStore()
      await store.loadExpenses()

      expect(store.expenses).toEqual(mockExpenses)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('正常系: 0件の場合空配列', async () => {
      vi.mocked(fetchExpenses).mockResolvedValue({
        success: true,
        data: [],
      })

      const store = useExpenseStore()
      await store.loadExpenses()

      expect(store.expenses).toEqual([])
      expect(store.error).toBeNull()
    })

    it('異常系: APIエラー', async () => {
      vi.mocked(fetchExpenses).mockResolvedValue({
        success: false,
        error: 'サーバーエラー',
      })

      const store = useExpenseStore()
      await store.loadExpenses()

      expect(store.error).toBe('サーバーエラー')
    })
  })

  describe('addExpense', () => {
    it('正常系: 経費作成成功', async () => {
      vi.mocked(createExpense).mockResolvedValue({ success: true, data: { id: 'new-id' } as any })
      vi.mocked(fetchExpenses).mockResolvedValue({ success: true, data: [] })

      const store = useExpenseStore()
      const result = await store.addExpense({
        expenseDate: '2026-04-20',
        categoryId: 'c1',
        amount: 1500,
        description: 'テスト経費',
      })

      expect(result).toBe(true)
      expect(createExpense).toHaveBeenCalledTimes(1)
    })

    it('異常系: 経費作成失敗', async () => {
      vi.mocked(createExpense).mockResolvedValue({ success: false, error: 'バリデーションエラー' })

      const store = useExpenseStore()
      const result = await store.addExpense({
        expenseDate: '2026-04-20',
        categoryId: 'c1',
        amount: 0,
        description: '',
      })

      expect(result).toBe(false)
      expect(store.error).toBe('バリデーションエラー')
    })
  })

  describe('removeExpense', () => {
    it('正常系: 経費削除成功', async () => {
      vi.mocked(deleteExpense).mockResolvedValue({ success: true })
      vi.mocked(fetchExpenses).mockResolvedValue({ success: true, data: [] })

      const store = useExpenseStore()
      const result = await store.removeExpense('expense-id')

      expect(result).toBe(true)
      expect(deleteExpense).toHaveBeenCalledWith('expense-id')
    })
  })

  describe('submit', () => {
    it('正常系: 申請成功', async () => {
      vi.mocked(submitExpense).mockResolvedValue({ success: true, data: { id: '1', status: 'pending_manager' } as any })
      vi.mocked(fetchExpenses).mockResolvedValue({ success: true, data: [] })

      const store = useExpenseStore()
      const result = await store.submit('expense-id')

      expect(result).toBe(true)
      expect(submitExpense).toHaveBeenCalledWith('expense-id')
    })

    it('異常系: 申請失敗', async () => {
      vi.mocked(submitExpense).mockResolvedValue({ success: false, error: '申請できない状態です' })

      const store = useExpenseStore()
      const result = await store.submit('expense-id')

      expect(result).toBe(false)
      expect(store.error).toBe('申請できない状態です')
    })
  })

  describe('loadCategories', () => {
    it('正常系: カテゴリ取得', async () => {
      const mockCategories = [
        { id: 'c1', code: 'transportation', name: '交通費', sortOrder: 1, isActive: true },
        { id: 'c2', code: 'entertainment', name: '交際費', sortOrder: 2, isActive: true },
      ]
      vi.mocked(fetchCategories).mockResolvedValue({ success: true, data: mockCategories })

      const store = useExpenseStore()
      await store.loadCategories()

      expect(store.categories).toEqual(mockCategories)
    })

    it('異常系: カテゴリ取得失敗（エラーは非表示）', async () => {
      vi.mocked(fetchCategories).mockRejectedValue(new Error('Network error'))

      const store = useExpenseStore()
      await store.loadCategories()

      // カテゴリ取得失敗は致命的でないためエラーは設定されない
      expect(store.error).toBeNull()
    })
  })
})
