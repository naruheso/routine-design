<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useExpenseStore } from '@/stores/expense'
import AppLayout from '@/components/common/AppLayout.vue'

const store = useExpenseStore()
const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const expenseId = computed(() => route.params.id as string)

const form = ref({ expenseDate: '', categoryId: '', amount: 0, description: '' })
const formError = ref<string | null>(null)

onMounted(async () => {
  await store.loadCategories()
  if (isEdit.value) {
    await store.loadExpenseDetail(expenseId.value)
    if (store.currentExpense) {
      const e = store.currentExpense.expense
      form.value = { expenseDate: e.expenseDate, categoryId: e.categoryId, amount: e.amount, description: e.description }
    }
  }
})

async function handleSave() {
  formError.value = null
  if (!form.value.expenseDate) { formError.value = '日付を入力してください'; return }
  if (!form.value.categoryId) { formError.value = '勘定科目を選択してください'; return }
  if (form.value.amount < 1) { formError.value = '金額は1円以上で入力してください'; return }
  if (form.value.amount > 999999) { formError.value = '金額は999,999円以下で入力してください'; return }
  if (!form.value.description) { formError.value = '摘要を入力してください'; return }

  let ok: boolean
  if (isEdit.value) {
    ok = await store.editExpense(expenseId.value, form.value)
  } else {
    ok = await store.addExpense(form.value)
  }
  if (ok) router.push('/')
  else formError.value = store.error
}

async function handleSubmit() {
  formError.value = null
  if (isEdit.value) {
    const ok = await store.editExpense(expenseId.value, form.value)
    if (!ok) { formError.value = store.error; return }
    const submitOk = await store.submit(expenseId.value)
    if (submitOk) router.push('/')
    else formError.value = store.error
  } else {
    const ok = await store.addExpense(form.value)
    if (!ok) { formError.value = store.error; return }
    // 作成後に最新の経費IDで申請
    if (store.expenses.length > 0) {
      const submitOk = await store.submit(store.expenses[0]!.id)
      if (submitOk) router.push('/')
      else formError.value = store.error
    }
  }
}
</script>

<template>
  <AppLayout>
    <div class="form-view">
      <h2>{{ isEdit ? '経費編集' : '経費入力' }}</h2>
      <div v-if="formError || store.error" class="err">{{ formError || store.error }}</div>
      <form class="expense-form" @submit.prevent="handleSave">
        <div class="fg">
          <label for="expenseDate">日付 <span class="req">*</span></label>
          <input id="expenseDate" v-model="form.expenseDate" type="date" :max="new Date().toISOString().slice(0,10)" required />
        </div>
        <div class="fg">
          <label for="categoryId">勘定科目 <span class="req">*</span></label>
          <select id="categoryId" v-model="form.categoryId" required>
            <option value="" disabled>選択してください</option>
            <option v-for="c in store.categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="fg">
          <label for="amount">金額 <span class="req">*</span></label>
          <input id="amount" v-model.number="form.amount" type="number" min="1" max="999999" placeholder="例: 1500" required />
        </div>
        <div class="fg">
          <label for="description">摘要 <span class="req">*</span></label>
          <textarea id="description" v-model="form.description" maxlength="500" rows="3" placeholder="経費の用途や詳細を入力" required></textarea>
          <span class="char-count">{{ form.description.length }}/500</span>
        </div>
        <div class="form-actions">
          <button type="button" class="btn-c" @click="router.push('/')">キャンセル</button>
          <button type="submit" class="btn-s" :disabled="store.loading">下書き保存</button>
          <button type="button" class="btn-p" :disabled="store.loading" @click="handleSubmit">申請する</button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>

<style scoped>
.form-view{max-width:600px}
.form-view h2{margin-bottom:var(--spacing-lg)}
.expense-form{display:flex;flex-direction:column;gap:var(--spacing-md);background:var(--color-surface-card);padding:var(--spacing-xl);border-radius:var(--radius-card);box-shadow:var(--shadow-card)}
.fg{display:flex;flex-direction:column;gap:var(--spacing-xs)}
.fg label{font-size:.875rem;font-weight:500;color:var(--color-text-heading)}
.req{color:var(--color-danger)}
.fg input,.fg select,.fg textarea{border:1px solid var(--color-input-border);border-radius:var(--radius-input);padding:var(--spacing-sm) var(--spacing-md);font-size:1rem}
.fg input:focus,.fg select:focus,.fg textarea:focus{outline:none;border-color:var(--color-primary);box-shadow:0 0 0 2px rgba(40,100,240,.15)}
.char-count{font-size:.75rem;color:var(--color-text-muted);text-align:right}
.form-actions{display:flex;gap:var(--spacing-sm);justify-content:flex-end;margin-top:var(--spacing-md)}
.btn-c{background:var(--color-bg-light);border:1px solid var(--color-border);border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer}
.btn-s{background:#fff;border:2px solid var(--color-primary);color:var(--color-primary);border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer;font-weight:500}
.btn-p{background:var(--color-primary);color:#fff;border:none;border-radius:var(--radius-button);padding:var(--spacing-sm) var(--spacing-lg);cursor:pointer;font-weight:500}
.btn-p:hover{background:var(--color-primary-hover)}
.err{background:#fef2f2;border:1px solid var(--color-danger);color:var(--color-danger);padding:var(--spacing-sm) var(--spacing-md);border-radius:var(--radius-input);margin-bottom:var(--spacing-md)}
</style>
