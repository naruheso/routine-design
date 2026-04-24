import { describe, it, expect, vi, beforeEach } from 'vitest'

// fetchモック
const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

// localStorageモック
const mockLocalStorage: Record<string, string> = {}
vi.stubGlobal('localStorage', {
  getItem: vi.fn((key: string) => mockLocalStorage[key] ?? null),
  setItem: vi.fn((key: string, value: string) => { mockLocalStorage[key] = value }),
  removeItem: vi.fn((key: string) => { delete mockLocalStorage[key] }),
  clear: vi.fn(() => { Object.keys(mockLocalStorage).forEach(k => delete mockLocalStorage[k]) }),
})

// window.location モック
const mockLocation = { href: '' }
vi.stubGlobal('location', mockLocation)

// apiClient をモック適用後にインポート
import { apiClient } from '@/api/client'

describe('apiClient', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    Object.keys(mockLocalStorage).forEach(k => delete mockLocalStorage[k])
    mockLocation.href = ''
  })

  describe('get', () => {
    it('正常系: GETリクエスト', async () => {
      const mockResponse = { success: true, data: { id: '1' } }
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      })

      const result = await apiClient.get('/expenses')

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/expenses',
        expect.objectContaining({ method: 'GET' })
      )
      expect(result.success).toBe(true)
    })

    it('正常系: JWTトークンが自動付与される', async () => {
      mockLocalStorage['token'] = 'test-jwt-token'
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true, data: [] }),
      })

      await apiClient.get('/expenses')

      const callArgs = mockFetch.mock.calls[0]
      expect(callArgs[1].headers['Authorization']).toBe('Bearer test-jwt-token')
    })
  })

  describe('post', () => {
    it('正常系: POSTリクエスト（ボディ付き）', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ success: true, data: { token: 'new-token' } }),
      })

      const result = await apiClient.post('/auth/login', { email: 'test@example.com', password: 'password123' })

      expect(mockFetch).toHaveBeenCalledWith(
        '/api/auth/login',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'test@example.com', password: 'password123' }),
        })
      )
      expect(result.success).toBe(true)
    })
  })

  describe('エラーハンドリング', () => {
    it('異常系: 401で認証エラー処理', async () => {
      mockLocalStorage['token'] = 'expired-token'
      mockFetch.mockResolvedValue({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ error: 'unauthorized' }),
      })

      const result = await apiClient.get('/expenses')

      expect(result.success).toBe(false)
    })

    it('異常系: 500エラー', async () => {
      mockFetch.mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.resolve({ error: 'Internal Server Error' }),
      })

      const result = await apiClient.get('/expenses')

      expect(result.success).toBe(false)
    })

    it('異常系: ネットワークエラー（fetch例外）', async () => {
      mockFetch.mockRejectedValue(new Error('Failed to fetch'))

      const result = await apiClient.get('/expenses')

      expect(result.success).toBe(false)
      expect(result.error).toBeDefined()
    })
  })
})
