import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'expenses',
      component: () => import('@/views/ExpenseListView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/expenses/new',
      name: 'expense-create',
      component: () => import('@/views/ExpenseFormView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/expenses/:id',
      name: 'expense-detail',
      component: () => import('@/views/ExpenseDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/expenses/:id/edit',
      name: 'expense-edit',
      component: () => import('@/views/ExpenseFormView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/approvals',
      name: 'approvals',
      component: () => import('@/views/ApprovalListView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import('@/views/ReportView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// 認証ガード
router.beforeEach((to) => {
  const authStore = useAuthStore()
  authStore.initialize()

  if (to.meta.requiresAuth !== false && !authStore.isAuthenticated) {
    return { name: 'login' }
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return { name: 'expenses' }
  }
})

export default router
