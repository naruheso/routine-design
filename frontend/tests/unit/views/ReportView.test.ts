import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import ReportView from '@/views/ReportView.vue'
import { useReportStore } from '@/stores/report'
import { useAuthStore } from '@/stores/auth'

// AppLayout をスタブ化
const AppLayoutStub = {
  template: '<div class="app-layout-stub"><slot /></div>'
}

describe('ReportView.vue', () => {
  let pinia: ReturnType<typeof createTestingPinia>

  beforeEach(() => {
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        auth: {
          user: { name: 'テストユーザー', roles: ['expense_admin'] }
        }
      }
    })
    vi.clearAllMocks()
  })

  it('正常系: 初期表示時にデータをロードする', () => {
    mount(ReportView, {
      global: {
        plugins: [pinia],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    const store = useReportStore()
    expect(store.loadAll).toHaveBeenCalled()
  })

  it('正常系: ローディング中にメッセージを表示する', () => {
    const store = useReportStore()
    // @ts-ignore: mock state
    store.loading = true

    const wrapper = mount(ReportView, {
      global: {
        plugins: [pinia],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    expect(wrapper.text()).toContain('読み込み中...')
    expect(wrapper.find('table').exists()).toBe(false)
  })

  it('正常系: データが存在する場合、テーブルに表示する', async () => {
    const store = useReportStore()
    // @ts-ignore: mock state
    store.loading = false
    // @ts-ignore: mock state
    store.employeeReport = {
      summary: { year: 2026, month: 4, totalCount: 2, grandTotal: 3000 },
      employees: [
        { userId: 'u1', userName: '山田太郎', count: 2, totalAmount: 3000 }
      ]
    }

    const wrapper = mount(ReportView, {
      global: {
        plugins: [pinia],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    expect(wrapper.find('#report-by-employee').exists()).toBe(true)
    expect(wrapper.text()).toContain('山田太郎')
    expect(wrapper.text()).toContain('¥3,000')
  })

  it('正常系: データが空の場合、メッセージを表示する', async () => {
    const store = useReportStore()
    // @ts-ignore: mock state
    store.loading = false
    // @ts-ignore: mock state
    store.employeeReport = { summary: { totalCount: 0, grandTotal: 0 }, employees: [] }

    const wrapper = mount(ReportView, {
      global: {
        plugins: [pinia],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    expect(wrapper.text()).toContain('該当するデータはありません')
  })

  it('正常系: 年月を変更すると再ロードされる', async () => {
    const wrapper = mount(ReportView, {
      global: {
        plugins: [pinia],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    const store = useReportStore()
    const yearSelect = wrapper.find('#report-year')
    const monthSelect = wrapper.find('#report-month')

    await yearSelect.setValue(2025)
    await monthSelect.setValue(12)

    // 初期化時 + 変更2回 ( setValueごとにwatchが反応 )
    expect(store.loadAll).toHaveBeenCalledWith(2025, 12)
  })
})
