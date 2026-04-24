// API レスポンス共通型
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  count?: number
}

export interface PaginatedResponse<T> extends ApiResponse<T> {
  pagination?: {
    currentPage: number
    totalPages: number
    totalCount: number
    limit: number
    hasNext: boolean
    hasPrev: boolean
  }
}
