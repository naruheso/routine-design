import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { LoginRequest, LoginResponse } from '@/types/user'

export function login(req: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  return apiClient.post<LoginResponse>('/auth/login', req)
}
