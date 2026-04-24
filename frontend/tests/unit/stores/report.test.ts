import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useReportStore } from '@/stores/report'

// APIモック
vi.mock('@/api/reports', () => ({
  fetchEmployeeReport: vi.fn(),
  fetchCategoryReport: vi.fn(),
}))

import { fetchEmployeeReport, fetchCategoryReport } from '@/api/reports'

describe('useReportStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('loadAll', () => {
    it('正常系: 社員別・勘定科目別の両方を取得', async () => {
      const mockEmpReport = {
        summary: { year: 2026, month: 4, totalCount: 3, grandTotal: 15000 },
        employees: [{ userId: 'u1', userName: '太郎', count: 3, totalAmount: 15000 }],
      }
      const mockCatReport = {
        summary: { year: 2026, month: 4, totalCount: 3, grandTotal: 15000 },
        categories: [{ categoryId: 'c1', categoryName: '交通費', count: 3, totalAmount: 15000, percentage: 100.0 }],
      }

      vi.mocked(fetchEmployeeReport).mockResolvedValue({ success: true, data: mockEmpReport })
      vi.mocked(fetchCategoryReport).mockResolvedValue({ success: true, data: mockCatReport })

      const store = useReportStore()
      await store.loadAll(2026, 4)

      expect(store.employeeReport).toEqual(mockEmpReport)
      expect(store.categoryReport).toEqual(mockCatReport)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(fetchEmployeeReport).toHaveBeenCalledWith(2026, 4)
      expect(fetchCategoryReport).toHaveBeenCalledWith(2026, 4)
    })

    it('正常系: 0件（データなし）', async () => {
      vi.mocked(fetchEmployeeReport).mockResolvedValue({
        success: true,
        data: { summary: { year: 2026, month: 4, totalCount: 0, grandTotal: 0 }, employees: [] },
      })
      vi.mocked(fetchCategoryReport).mockResolvedValue({
        success: true,
        data: { summary: { year: 2026, month: 4, totalCount: 0, grandTotal: 0 }, categories: [] },
      })

      const store = useReportStore()
      await store.loadAll(2026, 4)

      expect(store.employeeReport?.employees).toEqual([])
      expect(store.categoryReport?.categories).toEqual([])
      expect(store.error).toBeNull()
    })

    it('異常系: APIエラー', async () => {
      vi.mocked(fetchEmployeeReport).mockResolvedValue({ success: false, error: 'サーバーエラー' })
      vi.mocked(fetchCategoryReport).mockResolvedValue({ success: true, data: { summary: { year: 2026, month: 4, totalCount: 0, grandTotal: 0 }, categories: [] } })

      const store = useReportStore()
      await store.loadAll(2026, 4)

      expect(store.error).toBe('サーバーエラー')
    })

    it('異常系: ネットワークエラー', async () => {
      vi.mocked(fetchEmployeeReport).mockRejectedValue(new Error('Network error'))
      vi.mocked(fetchCategoryReport).mockRejectedValue(new Error('Network error'))

      const store = useReportStore()
      await store.loadAll(2026, 4)

      expect(store.error).toBe('Network error')
    })
  })

  describe('loadEmployeeReport', () => {
    it('正常系: 社員別レポート単体取得', async () => {
      const mockReport = {
        summary: { year: 2026, month: 1, totalCount: 5, grandTotal: 50000 },
        employees: [
          { userId: 'u1', userName: '太郎', count: 3, totalAmount: 30000 },
          { userId: 'u2', userName: '花子', count: 2, totalAmount: 20000 },
        ],
      }
      vi.mocked(fetchEmployeeReport).mockResolvedValue({ success: true, data: mockReport })

      const store = useReportStore()
      await store.loadEmployeeReport(2026, 1)

      expect(store.employeeReport?.employees).toHaveLength(2)
      expect(store.employeeReport?.summary.grandTotal).toBe(50000)
    })
  })
})
