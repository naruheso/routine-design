import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// APIモック
vi.mock('@/api/auth', () => ({
  login: vi.fn(),
}))

import { login } from '@/api/auth'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('初期状態', () => {
    it('正常系: 初期値が正しい', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('login', () => {
    it('正常系: ログイン成功', async () => {
      const mockUser = { id: '1', name: '太郎', email: 'test@example.com', roles: ['applicant'], createdAt: '', updatedAt: '' }
      vi.mocked(login).mockResolvedValue({
        success: true,
        data: { token: 'jwt-token', user: mockUser },
      })

      const store = useAuthStore()
      const result = await store.login('test@example.com', 'password123')

      expect(result).toBe(true)
      expect(store.isAuthenticated).toBe(true)
      expect(store.user).toEqual(mockUser)
      expect(store.token).toBe('jwt-token')
      expect(store.error).toBeNull()
      expect(localStorage.getItem('token')).toBe('jwt-token')
    })

    it('異常系: ログイン失敗（API エラー）', async () => {
      vi.mocked(login).mockResolvedValue({
        success: false,
        error: 'メールアドレスまたはパスワードが正しくありません',
      })

      const store = useAuthStore()
      const result = await store.login('wrong@example.com', 'wrong')

      expect(result).toBe(false)
      expect(store.isAuthenticated).toBe(false)
      expect(store.error).toBe('メールアドレスまたはパスワードが正しくありません')
    })

    it('異常系: ネットワークエラー', async () => {
      vi.mocked(login).mockRejectedValue(new Error('Network error'))

      const store = useAuthStore()
      const result = await store.login('test@example.com', 'password123')

      expect(result).toBe(false)
      expect(store.error).toBe('Network error')
    })
  })

  describe('logout', () => {
    it('正常系: ログアウトで状態クリア', async () => {
      // Arrange: ログイン状態にする
      vi.mocked(login).mockResolvedValue({
        success: true,
        data: { token: 'jwt-token', user: { id: '1', name: '太郎', email: 'test@example.com', roles: ['applicant'], createdAt: '', updatedAt: '' } },
      })

      const store = useAuthStore()
      await store.login('test@example.com', 'password123')
      expect(store.isAuthenticated).toBe(true)

      // Act
      store.logout()

      // Assert
      expect(store.isAuthenticated).toBe(false)
      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(localStorage.getItem('token')).toBeNull()
    })
  })

  describe('initialize', () => {
    it('正常系: localStorageからセッション復元', () => {
      const mockUser = { id: '1', name: '太郎', email: 'test@example.com', roles: ['applicant'] }
      localStorage.setItem('token', 'stored-token')
      localStorage.setItem('user', JSON.stringify(mockUser))

      const store = useAuthStore()
      store.initialize()

      expect(store.isAuthenticated).toBe(true)
      expect(store.token).toBe('stored-token')
      expect(store.user?.name).toBe('太郎')
    })

    it('異常系: localStorageのユーザー情報が壊れている', () => {
      localStorage.setItem('token', 'stored-token')
      localStorage.setItem('user', 'invalid-json')

      const store = useAuthStore()
      store.initialize()

      // パースエラー時はログアウト状態にリセット
      expect(store.isAuthenticated).toBe(false)
      expect(store.token).toBeNull()
    })

    it('境界値: localStorageが空', () => {
      const store = useAuthStore()
      store.initialize()

      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('hasRole', () => {
    it('正常系: ロール判定', async () => {
      vi.mocked(login).mockResolvedValue({
        success: true,
        data: { token: 'jwt-token', user: { id: '1', name: '太郎', email: 'test@example.com', roles: ['applicant', 'manager'], createdAt: '', updatedAt: '' } },
      })

      const store = useAuthStore()
      await store.login('test@example.com', 'password123')

      expect(store.hasRole('applicant')).toBe(true)
      expect(store.hasRole('manager')).toBe(true)
      expect(store.hasRole('finance_director')).toBe(false)
    })
  })

  describe('isApprover', () => {
    it('正常系: 承認者ロール判定', async () => {
      vi.mocked(login).mockResolvedValue({
        success: true,
        data: { token: 'jwt-token', user: { id: '1', name: '太郎', email: 'test@example.com', roles: ['manager'], createdAt: '', updatedAt: '' } },
      })

      const store = useAuthStore()
      await store.login('test@example.com', 'password123')

      expect(store.isApprover).toBe(true)
    })

    it('正常系: 申請者のみは承認者ではない', async () => {
      vi.mocked(login).mockResolvedValue({
        success: true,
        data: { token: 'jwt-token', user: { id: '1', name: '太郎', email: 'test@example.com', roles: ['applicant'], createdAt: '', updatedAt: '' } },
      })

      const store = useAuthStore()
      await store.login('test@example.com', 'password123')

      expect(store.isApprover).toBe(false)
    })
  })
})
