export interface Expense {
  id: string
  userId: string
  userName?: string
  expenseDate: string
  categoryId: string
  categoryName?: string
  amount: number
  description: string
  receiptImagePath: string | null
  status: string
  submittedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface CreateExpenseRequest {
  expenseDate: string
  categoryId: string
  amount: number
  description: string
}

export interface UpdateExpenseRequest {
  expenseDate: string
  categoryId: string
  amount: number
  description: string
}

export interface ApprovalHistory {
  id: string
  expenseId: string
  action: string
  actorId: string
  actorName?: string
  actorRole: string
  comment: string | null
  createdAt: string
}

export interface ExpenseDetail {
  expense: Expense
  histories: ApprovalHistory[]
}

export interface Category {
  id: string
  code: string
  name: string
  sortOrder: number
  isActive: boolean
}

export interface ApprovalRequest {
  comment?: string
}

// ステータス定数
export const STATUSES = {
  DRAFT: 'draft',
  PENDING_MANAGER: 'pending_manager',
  PENDING_EXPENSE_ADMIN: 'pending_expense_admin',
  PENDING_FINANCE_DIRECTOR: 'pending_finance_director',
  RETURNED: 'returned',
  REJECTED: 'rejected',
  APPROVED: 'approved',
} as const

// ステータス表示名
export const STATUS_DISPLAY_NAMES: Record<string, string> = {
  draft: '下書き',
  pending_manager: '上長承認待ち',
  pending_expense_admin: '経費担当承認待ち',
  pending_finance_director: '経理部長承認待ち',
  returned: '差し戻し',
  rejected: '否認',
  approved: '承認完了',
}

// ステータスの色（PrimeVue severity）
export const STATUS_SEVERITY: Record<string, string> = {
  draft: 'secondary',
  pending_manager: 'warn',
  pending_expense_admin: 'warn',
  pending_finance_director: 'warn',
  returned: 'danger',
  rejected: 'danger',
  approved: 'success',
}
