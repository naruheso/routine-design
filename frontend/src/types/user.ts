export interface User {
  id: string
  name: string
  email: string
  roles: string[]
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// ロール定数
export const ROLES = {
  APPLICANT: 'applicant',
  MANAGER: 'manager',
  EXPENSE_ADMIN: 'expense_admin',
  FINANCE_DIRECTOR: 'finance_director',
} as const

export type UserRole = (typeof ROLES)[keyof typeof ROLES]

// ロール表示名
export const ROLE_DISPLAY_NAMES: Record<string, string> = {
  applicant: '申請者（社員）',
  manager: '上長',
  expense_admin: '経費担当',
  finance_director: '経理部長',
}
