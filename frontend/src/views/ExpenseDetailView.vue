<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useExpenseStore } from '@/stores/expense'
import AppLayout from '@/components/common/AppLayout.vue'
import { STATUS_DISPLAY_NAMES, STATUS_SEVERITY, STATUSES } from '@/types/expense'
import { ROLE_DISPLAY_NAMES } from '@/types/user'

const store = useExpenseStore()
const route = useRoute()
const router = useRouter()
const id = route.params.id as string

onMounted(() => { store.loadExpenseDetail(id) })

function formatAmount(n: number) { return n.toLocaleString('ja-JP') + '円' }
function badgeClass(s: string) { return `badge--${STATUS_SEVERITY[s] ?? 'secondary'}` }
function actionName(a: string) {
  const m: Record<string,string> = { submit:'申請', approve:'承認', return:'差し戻し', reject:'否認' }
  return m[a] ?? a
}

async function handleSubmit() {
  const ok = await store.submit(id)
  if (ok) router.push('/')
}
</script>

<template>
  <AppLayout>
    <div v-if="store.loading" class="empty">読み込み中...</div>
    <div v-else-if="!store.currentExpense" class="empty">データが見つかりません</div>
    <div v-else class="detail">
      <div class="detail__header">
        <h2>経費詳細</h2>
        <button class="btn-back" @click="router.push('/')">← 一覧に戻る</button>
      </div>
      <div class="card">
        <div class="info-grid">
          <div class="info"><span class="info__label">ステータス</span><span class="badge" :class="badgeClass(store.currentExpense.expense.status)">{{ STATUS_DISPLAY_NAMES[store.currentExpense.expense.status] }}</span></div>
          <div class="info"><span class="info__label">日付</span><span>{{ store.currentExpense.expense.expenseDate }}</span></div>
          <div class="info"><span class="info__label">勘定科目</span><span>{{ store.currentExpense.expense.categoryName }}</span></div>
          <div class="info"><span class="info__label">金額</span><span>{{ formatAmount(store.currentExpense.expense.amount) }}</span></div>
          <div class="info"><span class="info__label">摘要</span><span>{{ store.currentExpense.expense.description }}</span></div>
          <div class="info"><span class="info__label">申請者</span><span>{{ store.currentExpense.expense.userName }}</span></div>
        </div>
        <div v-if="store.currentExpense.expense.status===STATUSES.DRAFT||store.currentExpense.expense.status===STATUSES.RETURNED" class="detail__actions">
          <button class="btn-edit" @click="router.push(`/expenses/${id}/edit`)">編集</button>
          <button v-if="store.currentExpense.expense.status===STATUSES.DRAFT||store.currentExpense.expense.status===STATUSES.RETURNED" class="btn-submit" @click="handleSubmit">申請する</button>
        </div>
      </div>
      <div v-if="store.currentExpense.histories?.length" class="card history">
        <h3>承認履歴</h3>
        <div v-for="h in store.currentExpense.histories" :key="h.id" class="history__item">
          <div class="history__meta">
            <strong>{{ actionName(h.action) }}</strong>
            <span>{{ h.actorName }} ({{ ROLE_DISPLAY_NAMES[h.actorRole] ?? h.actorRole }})</span>
            <span class="history__time">{{ new Date(h.createdAt).toLocaleString('ja-JP') }}</span>
          </div>
          <div v-if="h.comment" class="history__comment">{{ h.comment }}</div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
.detail__header{display:flex;justify-content:space-between;align-items:center;margin-bottom:var(--spacing-lg)}
.btn-back{background:transparent;border:none;color:var(--color-primary);cursor:pointer;font-size:.875rem}
.card{background:var(--color-surface-card);border-radius:var(--radius-card);box-shadow:var(--shadow-card);padding:var(--spacing-xl);margin-bottom:var(--spacing-lg)}
.info-grid{display:grid;grid-template-columns:1fr 1fr;gap:var(--spacing-md)}
.info{display:flex;flex-direction:column;gap:var(--spacing-xs)}
.info__label{font-size:.75rem;color:var(--color-text-muted);font-weight:500}
.badge{display:inline-block;padding:2px var(--spacing-sm);border-radius:12px;font-size:.75rem;font-weight:500;width:fit-content}
.badge--secondary{background:#f0eded;color:#595959}.badge--warn{background:#fff8e1;color:#b8860b}
.badge--danger{background:#fef2f2;color:var(--color-danger)}.badge--success{background:#e8f5e9;color:var(--color-success)}
.detail__actions{display:flex;gap:var(--spacing-sm);margin-top:var(--spacing-lg);padding-top:var(--spacing-md);border-top:1px solid var(--color-border)}
.btn-edit{background:#fff;border:2px solid var(--color-primary);color:var(--color-primary);border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer;font-weight:500}
.btn-submit{background:var(--color-primary);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer;font-weight:500}
.history h3{margin-bottom:var(--spacing-md)}
.history__item{border-left:3px solid var(--color-primary);padding:var(--spacing-sm) var(--spacing-md);margin-bottom:var(--spacing-sm);background:var(--color-bg-light);border-radius:0 var(--radius-input) var(--radius-input) 0}
.history__meta{display:flex;gap:var(--spacing-sm);align-items:center;font-size:.875rem;flex-wrap:wrap}
.history__time{color:var(--color-text-muted);font-size:.75rem;margin-left:auto}
.history__comment{margin-top:var(--spacing-xs);font-size:.875rem;color:var(--color-text-body)}
.empty{text-align:center;padding:var(--spacing-xxl);color:var(--color-text-muted)}
</style>
