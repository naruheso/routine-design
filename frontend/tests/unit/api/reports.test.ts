import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchEmployeeReport, fetchCategoryReport } from '@/api/reports'

vi.mock('@/api/client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '@/api/client'

describe('reports API', () => {
  beforeEach(() => { vi.clearAllMocks() })

  describe('fetchEmployeeReport()', () => {
    it('正常系: 社員別レポートを取得する', async () => {
      const mockResponse = {
        summary: { year: 2026, month: 4, totalCount: 5, grandTotal: 50000 },
        employees: [
          { userId: 'user-1', userName: '田中太郎', count: 3, totalAmount: 30000 },
          { userId: 'user-2', userName: '佐藤花子', count: 2, totalAmount: 20000 },
        ],
      }
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: mockResponse })

      const result = await fetchEmployeeReport(2026, 4)

      expect(apiClient.get).toHaveBeenCalledWith('/reports/by-employee?year=2026&month=4')
      expect(result.success).toBe(true)
      expect(result.data?.employees).toHaveLength(2)
      expect(result.data?.summary.grandTotal).toBe(50000)
    })

    it('正常系: データが0件の場合', async () => {
      const mockResponse = {
        summary: { year: 2026, month: 1, totalCount: 0, grandTotal: 0 },
        employees: [],
      }
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: mockResponse })

      const result = await fetchEmployeeReport(2026, 1)
      expect(result.data?.employees).toEqual([])
      expect(result.data?.summary.grandTotal).toBe(0)
    })

    it('異常系: 取得失敗', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: false, error: '権限がありません' })
      const result = await fetchEmployeeReport(2026, 4)
      expect(result.success).toBe(false)
    })

    it('境界値: 年月の境界 (1月・12月)', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: { summary: { year: 2026, month: 12, totalCount: 0, grandTotal: 0 }, employees: [] } })
      await fetchEmployeeReport(2026, 12)
      expect(apiClient.get).toHaveBeenCalledWith('/reports/by-employee?year=2026&month=12')

      await fetchEmployeeReport(2026, 1)
      expect(apiClient.get).toHaveBeenCalledWith('/reports/by-employee?year=2026&month=1')
    })
  })

  describe('fetchCategoryReport()', () => {
    it('正常系: 勘定科目別レポートを取得する', async () => {
      const mockResponse = {
        summary: { year: 2026, month: 4, totalCount: 10, grandTotal: 100000 },
        categories: [
          { categoryId: 'cat-1', categoryName: '交通費', count: 5, totalAmount: 50000, percentage: 50.0 },
          { categoryId: 'cat-2', categoryName: '交際費', count: 5, totalAmount: 50000, percentage: 50.0 },
        ],
      }
      vi.mocked(apiClient.get).mockResolvedValue({ success: true, data: mockResponse })

      const result = await fetchCategoryReport(2026, 4)

      expect(apiClient.get).toHaveBeenCalledWith('/reports/by-category?year=2026&month=4')
      expect(result.success).toBe(true)
      expect(result.data?.categories).toHaveLength(2)
    })

    it('異常系: 取得失敗', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ success: false, error: 'サーバーエラー' })
      const result = await fetchCategoryReport(2026, 4)
      expect(result.success).toBe(false)
    })
  })
})
