<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useApprovalStore } from '@/stores/approval'
import AppLayout from '@/components/common/AppLayout.vue'
import { STATUS_DISPLAY_NAMES, STATUS_SEVERITY } from '@/types/expense'
import type { Expense } from '@/types/expense'

const store = useApprovalStore()
const router = useRouter()

const showDialog = ref(false)
const dialogAction = ref<'approve'|'return'|'reject'>('approve')
const dialogTarget = ref<Expense|null>(null)
const comment = ref('')
const dialogError = ref<string|null>(null)

onMounted(() => { store.loadPending() })

function formatAmount(n: number) { return n.toLocaleString('ja-JP') + '円' }
function badgeClass(s: string) { return `badge--${STATUS_SEVERITY[s] ?? 'secondary'}` }

function openDialog(action: 'approve'|'return'|'reject', expense: Expense) {
  dialogAction.value = action
  dialogTarget.value = expense
  comment.value = ''
  dialogError.value = null
  showDialog.value = true
}

const actionLabels: Record<string,string> = { approve:'承認', return:'差し戻し', reject:'否認' }

async function handleAction() {
  if (!dialogTarget.value) return
  dialogError.value = null
  const id = dialogTarget.value.id
  let ok = false
  if (dialogAction.value === 'approve') {
    ok = await store.approve(id, comment.value || undefined)
  } else if (dialogAction.value === 'return') {
    if (!comment.value) { dialogError.value = '差し戻しの理由をコメントに入力してください'; return }
    ok = await store.returnBack(id, comment.value)
  } else {
    if (!comment.value) { dialogError.value = '否認の理由をコメントに入力してください'; return }
    ok = await store.reject(id, comment.value)
  }
  if (ok) { showDialog.value = false }
  else { dialogError.value = store.error }
}
</script>

<template>
  <AppLayout>
    <div class="view">
      <h2>承認待ち一覧</h2>
      <div v-if="store.error" class="err">{{ store.error }}</div>
      <div v-if="store.loading" class="empty">読み込み中...</div>
      <div v-else-if="!store.pendingExpenses.length" class="empty">承認待ちの申請は0件です</div>
      <table v-else class="tbl">
        <thead><tr><th>申請者</th><th>日付</th><th>科目</th><th>金額</th><th>摘要</th><th>ステータス</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="e in store.pendingExpenses" :key="e.id" class="tbl__row" @click="router.push(`/expenses/${e.id}`)">
            <td>{{ e.userName }}</td>
            <td>{{ e.expenseDate }}</td>
            <td>{{ e.categoryName }}</td>
            <td class="r">{{ formatAmount(e.amount) }}</td>
            <td class="trunc">{{ e.description }}</td>
            <td><span class="badge" :class="badgeClass(e.status)">{{ STATUS_DISPLAY_NAMES[e.status] }}</span></td>
            <td @click.stop>
              <button class="sm sm-ok" @click="openDialog('approve',e)">承認</button>
              <button class="sm sm-rt" @click="openDialog('return',e)">差し戻し</button>
              <button class="sm sm-rj" @click="openDialog('reject',e)">否認</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="showDialog" class="overlay" @click.self="showDialog=false">
        <div class="dlg">
          <h3>{{ actionLabels[dialogAction] }}確認</h3>
          <p>{{ dialogTarget?.userName }}さんの申請（{{ formatAmount(dialogTarget?.amount??0) }}）を{{ actionLabels[dialogAction] }}しますか？</p>
          <div v-if="dialogError" class="err">{{ dialogError }}</div>
          <div class="fg">
            <label>コメント{{ dialogAction!=='approve'?' (必須)':' (任意)' }}</label>
            <textarea v-model="comment" rows="3" :required="dialogAction!=='approve'" :placeholder="dialogAction==='approve'?'コメントを入力（任意）':'理由を入力してください'"></textarea>
          </div>
          <div class="dlg__act">
            <button class="btn-c" @click="showDialog=false">キャンセル</button>
            <button :class="dialogAction==='approve'?'btn-ok':'btn-dng'" @click="handleAction">{{ actionLabels[dialogAction] }}する</button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
h2{margin-bottom:var(--spacing-lg)}
.tbl{width:100%;border-collapse:collapse;background:var(--color-surface-card);border-radius:var(--radius-card);overflow:hidden;box-shadow:var(--shadow-card)}
.tbl th{background:var(--color-bg-light);padding:var(--spacing-sm) var(--spacing-md);text-align:left;font-size:.8125rem;font-weight:600;border-bottom:2px solid var(--color-border)}
.tbl td{padding:var(--spacing-sm) var(--spacing-md);border-bottom:1px solid var(--color-border);font-size:.875rem}
.tbl__row{cursor:pointer}.tbl__row:hover{background:var(--color-bg-light)}
.r{text-align:right}.trunc{max-width:150px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.badge{display:inline-block;padding:2px var(--spacing-sm);border-radius:12px;font-size:.75rem;font-weight:500}
.badge--warn{background:#fff8e1;color:#b8860b}
.sm{padding:2px var(--spacing-sm);border-radius:var(--radius-input);font-size:.75rem;cursor:pointer;border:1px solid;margin-right:4px}
.sm-ok{background:transparent;border-color:var(--color-success);color:var(--color-success)}.sm-ok:hover{background:#e8f5e9}
.sm-rt{background:transparent;border-color:var(--color-warning);color:#b8860b}.sm-rt:hover{background:#fff8e1}
.sm-rj{background:transparent;border-color:var(--color-danger);color:var(--color-danger)}.sm-rj:hover{background:#fef2f2}
.overlay{position:fixed;inset:0;background:rgba(0,0,0,.4);display:flex;align-items:center;justify-content:center;z-index:1000}
.dlg{background:var(--color-surface-card);border-radius:var(--radius-dialog);padding:var(--spacing-xl);max-width:480px;width:90%}
.dlg h3{margin-bottom:var(--spacing-md)}
.fg{margin:var(--spacing-md) 0}.fg label{font-size:.875rem;font-weight:500;display:block;margin-bottom:var(--spacing-xs)}
.fg textarea{width:100%;border:1px solid var(--color-input-border);border-radius:var(--radius-input);padding:var(--spacing-sm);font-size:.875rem;resize:vertical}
.dlg__act{display:flex;gap:var(--spacing-sm);justify-content:flex-end;margin-top:var(--spacing-lg)}
.btn-c{background:var(--color-bg-light);border:1px solid var(--color-border);border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
.btn-ok{background:var(--color-success);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
.btn-dng{background:var(--color-danger);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
.err{background:#fef2f2;border:1px solid var(--color-danger);color:var(--color-danger);padding:var(--spacing-sm) var(--spacing-md);border-radius:var(--radius-input);margin-bottom:var(--spacing-md)}
.empty{text-align:center;padding:var(--spacing-xxl);color:var(--color-text-muted)}
</style>
