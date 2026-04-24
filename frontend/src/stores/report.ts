import { ref } from 'vue'
import { defineStore } from 'pinia'
import { fetchEmployeeReport, fetchCategoryReport } from '@/api/reports'
import type { EmployeeReportResponse, CategoryReportResponse } from '@/types/report'

export const useReportStore = defineStore('report', () => {
  // State
  const employeeReport = ref<EmployeeReportResponse | null>(null)
  const categoryReport = ref<CategoryReportResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function loadEmployeeReport(year: number, month: number) {
    loading.value = true
    error.value = null
    try {
      const result = await fetchEmployeeReport(year, month)
      if (result.success && result.data !== undefined) {
        employeeReport.value = result.data
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  async function loadCategoryReport(year: number, month: number) {
    loading.value = true
    error.value = null
    try {
      const result = await fetchCategoryReport(year, month)
      if (result.success && result.data !== undefined) {
        categoryReport.value = result.data
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  async function loadAll(year: number, month: number) {
    loading.value = true
    error.value = null
    try {
      const [empResult, catResult] = await Promise.all([
        fetchEmployeeReport(year, month),
        fetchCategoryReport(year, month),
      ])

      if (empResult.success && empResult.data !== undefined) {
        employeeReport.value = empResult.data
      }
      if (catResult.success && catResult.data !== undefined) {
        categoryReport.value = catResult.data
      }

      if (!empResult.success || !catResult.success) {
        error.value = empResult.error ?? catResult.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  return {
    employeeReport, categoryReport, loading, error,
    loadEmployeeReport, loadCategoryReport, loadAll,
  }
})
