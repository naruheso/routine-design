import { describe, it, expect, vi, beforeEach } from 'vitest'
import { login } from '@/api/auth'

// apiClient をモック
vi.mock('@/api/client', () => ({
  apiClient: {
    post: vi.fn(),
    get: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { apiClient } from '@/api/client'

describe('auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('login()', () => {
    it('正常系: 正しい認証情報でログイン成功', async () => {
      // Arrange
      const mockResponse = {
        success: true,
        data: {
          token: 'jwt-token-123',
          user: {
            id: 'user-1',
            name: '田中太郎',
            email: 'tanaka@example.com',
            roles: ['applicant'],
            createdAt: '2026-01-01T00:00:00Z',
            updatedAt: '2026-01-01T00:00:00Z',
          },
        },
      }
      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      // Act
      const result = await login({ email: 'tanaka@example.com', password: 'password123' })

      // Assert
      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', {
        email: 'tanaka@example.com',
        password: 'password123',
      })
      expect(result.success).toBe(true)
      expect(result.data?.token).toBe('jwt-token-123')
      expect(result.data?.user.name).toBe('田中太郎')
    })

    it('異常系: 認証情報が誤っている場合', async () => {
      // Arrange
      const mockResponse = {
        success: false,
        error: 'メールアドレスまたはパスワードが正しくありません',
      }
      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      // Act
      const result = await login({ email: 'wrong@example.com', password: 'wrong' })

      // Assert
      expect(result.success).toBe(false)
      expect(result.error).toBe('メールアドレスまたはパスワードが正しくありません')
    })

    it('境界値: 空文字のメールアドレスとパスワードでリクエスト送信', async () => {
      // Arrange
      vi.mocked(apiClient.post).mockResolvedValue({ success: false, error: 'バリデーションエラー' })

      // Act
      const result = await login({ email: '', password: '' })

      // Assert
      expect(apiClient.post).toHaveBeenCalledWith('/auth/login', { email: '', password: '' })
      expect(result.success).toBe(false)
    })
  })
})
