<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useExpenseStore } from '@/stores/expense'
import AppLayout from '@/components/common/AppLayout.vue'
import { STATUS_DISPLAY_NAMES, STATUS_SEVERITY, STATUSES } from '@/types/expense'

const store = useExpenseStore()
const router = useRouter()
const showDeleteConfirm = ref(false)
const deleteTargetId = ref('')

onMounted(() => { store.loadExpenses() })

function goToCreate() { router.push('/expenses/new') }
function goToDetail(id: string) { router.push(`/expenses/${id}`) }
function goToEdit(id: string) { router.push(`/expenses/${id}/edit`) }

function confirmDelete(id: string) {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (deleteTargetId.value) {
    await store.removeExpense(deleteTargetId.value)
    showDeleteConfirm.value = false
  }
}

function formatAmount(n: number) { return n.toLocaleString('ja-JP') + '円' }
function badgeClass(s: string) { return `badge--${STATUS_SEVERITY[s] ?? 'secondary'}` }
</script>

<template>
  <AppLayout>
    <div class="view">
      <div class="view__header">
        <h2>経費申請一覧</h2>
        <button class="btn-primary" @click="goToCreate">＋ 新規作成</button>
      </div>
      <div v-if="store.error" class="err">{{ store.error }}</div>
      <div v-if="store.loading" class="empty">読み込み中...</div>
      <div v-else-if="!store.expenses.length" class="empty">
        <p>申請は0件です</p>
        <button class="btn-primary" @click="goToCreate">最初の経費を入力する</button>
      </div>
      <table v-else class="tbl">
        <thead><tr><th>日付</th><th>勘定科目</th><th>金額</th><th>摘要</th><th>ステータス</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="e in store.expenses" :key="e.id" class="tbl__row" @click="goToDetail(e.id)">
            <td>{{ e.expenseDate }}</td>
            <td>{{ e.categoryName }}</td>
            <td class="r">{{ formatAmount(e.amount) }}</td>
            <td class="trunc">{{ e.description }}</td>
            <td><span class="badge" :class="badgeClass(e.status)">{{ STATUS_DISPLAY_NAMES[e.status] }}</span></td>
            <td @click.stop>
              <button v-if="e.status===STATUSES.DRAFT||e.status===STATUSES.RETURNED" class="sm sm-e" @click="goToEdit(e.id)">編集</button>
              <button v-if="e.status===STATUSES.DRAFT" class="sm sm-d" @click="confirmDelete(e.id)">削除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="showDeleteConfirm" class="overlay" @click.self="showDeleteConfirm=false">
        <div class="dlg"><h3>削除確認</h3><p>この経費データを削除しますか？</p>
          <div class="dlg__act"><button class="btn-c" @click="showDeleteConfirm=false">キャンセル</button><button class="btn-dng" @click="handleDelete">削除する</button></div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
.view__header{display:flex;justify-content:space-between;align-items:center;margin-bottom:var(--spacing-lg)}
.view__header h2{font-size:1.25rem}
.btn-primary{background:var(--color-primary);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);font-weight:500;cursor:pointer}
.btn-primary:hover{background:var(--color-primary-hover)}
.tbl{width:100%;border-collapse:collapse;background:var(--color-surface-card);border-radius:var(--radius-card);overflow:hidden;box-shadow:var(--shadow-card)}
.tbl th{background:var(--color-bg-light);padding:var(--spacing-sm) var(--spacing-md);text-align:left;font-size:.8125rem;font-weight:600;border-bottom:2px solid var(--color-border)}
.tbl td{padding:var(--spacing-sm) var(--spacing-md);border-bottom:1px solid var(--color-border);font-size:.875rem}
.tbl__row{cursor:pointer}.tbl__row:hover{background:var(--color-bg-light)}
.r{text-align:right}.trunc{max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.badge{display:inline-block;padding:2px var(--spacing-sm);border-radius:12px;font-size:.75rem;font-weight:500}
.badge--secondary{background:#f0eded;color:#595959}.badge--warn{background:#fff8e1;color:#b8860b}
.badge--danger{background:#fef2f2;color:var(--color-danger)}.badge--success{background:#e8f5e9;color:var(--color-success)}
.sm{padding:2px var(--spacing-sm);border-radius:var(--radius-input);font-size:.75rem;cursor:pointer;border:1px solid;margin-right:4px}
.sm-e{background:transparent;border-color:var(--color-primary);color:var(--color-primary)}.sm-e:hover{background:rgba(40,100,240,.08)}
.sm-d{background:transparent;border-color:var(--color-danger);color:var(--color-danger)}.sm-d:hover{background:#fef2f2}
.err{background:#fef2f2;border:1px solid var(--color-danger);color:var(--color-danger);padding:var(--spacing-sm) var(--spacing-md);border-radius:var(--radius-input);margin-bottom:var(--spacing-md)}
.empty{text-align:center;padding:var(--spacing-xxl);color:var(--color-text-muted)}
.overlay{position:fixed;inset:0;background:rgba(0,0,0,.4);display:flex;align-items:center;justify-content:center;z-index:1000}
.dlg{background:var(--color-surface-card);border-radius:var(--radius-dialog);padding:var(--spacing-xl);max-width:400px;width:90%}
.dlg h3{margin-bottom:var(--spacing-md)}.dlg__act{display:flex;gap:var(--spacing-sm);justify-content:flex-end;margin-top:var(--spacing-lg)}
.btn-c{background:var(--color-bg-light);border:1px solid var(--color-border);border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
.btn-dng{background:var(--color-danger);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
</style>
