import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { login as loginApi } from '@/api/auth'
import type { User } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const userName = computed(() => user.value?.name ?? '')
  const userRoles = computed(() => user.value?.roles ?? [])

  const hasRole = (role: string) => {
    return userRoles.value.includes(role)
  }

  const isApprover = computed(() => {
    return hasRole('manager') || hasRole('expense_admin') || hasRole('finance_director')
  })

  // Actions
  async function login(email: string, password: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const result = await loginApi({ email, password })
      if (result.success && result.data) {
        token.value = result.data.token
        user.value = result.data.user
        localStorage.setItem('token', result.data.token)
        localStorage.setItem('user', JSON.stringify(result.data.user))
        return true
      } else {
        error.value = result.error ?? 'ログインに失敗しました'
        return false
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 初期化: localStorageからセッション復元
  function initialize() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    if (storedToken && storedUser) {
      token.value = storedToken
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        logout()
      }
    }
  }

  return {
    user, token, loading, error,
    isAuthenticated, userName, userRoles, isApprover,
    hasRole, login, logout, initialize,
  }
})
