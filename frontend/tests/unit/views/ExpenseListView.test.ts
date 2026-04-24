import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createRouter, createMemoryHistory } from 'vue-router'
import ExpenseListView from '@/views/ExpenseListView.vue'
import { useExpenseStore } from '@/stores/expense'

// Routerモック
const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: { template: 'Home' } },
    { path: '/expenses/new', component: { template: 'New' } },
    { path: '/expenses/:id', component: { template: 'Detail' } },
  ],
})

const AppLayoutStub = {
  template: '<div class="app-layout-stub"><slot /></div>'
}

describe('ExpenseListView.vue', () => {
  let pinia: ReturnType<typeof createTestingPinia>

  beforeEach(() => {
    pinia = createTestingPinia({ createSpy: vi.fn })
    vi.clearAllMocks()
  })

  it('正常系: 初期表示時に経費一覧をロードする', () => {
    mount(ExpenseListView, {
      global: {
        plugins: [pinia, router],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    const store = useExpenseStore()
    expect(store.loadExpenses).toHaveBeenCalled()
  })

  it('正常系: 経費が存在する場合、テーブルに表示される', async () => {
    const store = useExpenseStore()
    // @ts-ignore: mock state
    store.expenses = [
      { id: '1', expenseDate: '2026-04-24', categoryName: '交通費', amount: 1500, description: '電車代', status: 'draft' }
    ]
    // @ts-ignore: mock state
    store.loading = false

    const wrapper = mount(ExpenseListView, {
      global: {
        plugins: [pinia, router],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.text()).toContain('交通費')
    expect(wrapper.text()).toContain('1,500円')
  })

  it('正常系: 新規作成ボタンクリックで遷移する', async () => {
    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = mount(ExpenseListView, {
      global: {
        plugins: [pinia, router],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    await wrapper.find('.btn-primary').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/expenses/new')
  })

  it('正常系: 行をクリックすると詳細画面へ遷移する', async () => {
    const store = useExpenseStore()
    // @ts-ignore: mock state
    store.expenses = [{ id: '99', expenseDate: '2026-04-24', categoryName: '交通費', amount: 1500, status: 'draft' }]
    // @ts-ignore: mock state
    store.loading = false

    const pushSpy = vi.spyOn(router, 'push')
    const wrapper = mount(ExpenseListView, {
      global: {
        plugins: [pinia, router],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    await wrapper.find('.tbl__row').trigger('click')
    expect(pushSpy).toHaveBeenCalledWith('/expenses/99')
  })

  it('正常系: 下書き状態のみ削除ボタンが表示される', async () => {
    const store = useExpenseStore()
    // @ts-ignore: mock state
    store.expenses = [
      { id: '1', status: 'draft', amount: 100, expenseDate: '2026-04-24', categoryName: '交通費' },
      { id: '2', status: 'pending_manager', amount: 200, expenseDate: '2026-04-24', categoryName: '交通費' }
    ]
    // @ts-ignore: mock state
    store.loading = false

    const wrapper = mount(ExpenseListView, {
      global: {
        plugins: [pinia, router],
        stubs: { AppLayout: AppLayoutStub }
      }
    })

    const rows = wrapper.findAll('.tbl__row')
    expect(rows[0].find('.sm-d').exists()).toBe(true)
    expect(rows[1].find('.sm-d').exists()).toBe(false)
  })
})
