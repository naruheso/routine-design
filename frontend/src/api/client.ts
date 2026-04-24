import type { ApiResponse } from '@/types/api'

const API_BASE = '/api'

// 共通リクエスト関数
async function request<T>(path: string, options?: RequestInit): Promise<ApiResponse<T>> {
  try {
    const token = localStorage.getItem('token')
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const res = await fetch(`${API_BASE}${path}`, {
      headers,
      ...options,
    })

    if (res.status === 401) {
      // 認証エラー時はトークン削除してログインへリダイレクト
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
      return { success: false, error: '認証が必要です' }
    }

    if (!res.ok) {
      const body = await res.json().catch(() => null)
      return { success: false, error: body?.error ?? `HTTP ${res.status}` }
    }

    return await res.json()
  } catch (e) {
    return { success: false, error: e instanceof Error ? e.message : '通信エラー' }
  }
}

// APIクライアント
export const apiClient = {
  get<T>(path: string): Promise<ApiResponse<T>> {
    return request<T>(path, { method: 'GET' })
  },

  post<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  put<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
    return request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  delete<T>(path: string): Promise<ApiResponse<T>> {
    return request<T>(path, { method: 'DELETE' })
  },
}
