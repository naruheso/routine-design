<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useReportStore } from '@/stores/report'
import AppLayout from '@/components/common/AppLayout.vue'

const store = useReportStore()
const now = new Date()
const selectedYear = ref(now.getFullYear())
const selectedMonth = ref(now.getMonth() + 1)

// 年の選択肢（2020年〜今年＋1年）
const years = Array.from({ length: now.getFullYear() - 2020 + 2 }, (_, i) => 2020 + i)
const months = Array.from({ length: 12 }, (_, i) => i + 1)

function formatAmount(n: number) { return '¥' + n.toLocaleString('ja-JP') }
function formatPct(n: number) { return n.toFixed(1) + '%' }

function loadData() {
  store.loadAll(selectedYear.value, selectedMonth.value)
}

onMounted(() => { loadData() })

watch([selectedYear, selectedMonth], () => { loadData() })
</script>

<template>
  <AppLayout>
    <div class="report-view">
      <div class="report-view__header">
        <h2>集計レポート</h2>
        <div class="year-month-selector">
          <label>対象年月:</label>
          <select id="report-year" v-model="selectedYear">
            <option v-for="y in years" :key="y" :value="y">{{ y }}年</option>
          </select>
          <select id="report-month" v-model="selectedMonth">
            <option v-for="m in months" :key="m" :value="m">{{ m }}月</option>
          </select>
        </div>
      </div>

      <div v-if="store.error" class="err">{{ store.error }}</div>
      <div v-if="store.loading" class="empty">読み込み中...</div>

      <template v-else>
        <!-- 社員別 月別合計 -->
        <section class="report-section">
          <h3>社員別 月別合計</h3>
          <div v-if="!store.employeeReport?.employees?.length" class="empty-section">
            該当するデータはありません
          </div>
          <table v-else class="tbl" id="report-by-employee">
            <thead>
              <tr><th>社員名</th><th class="r">件数</th><th class="r">合計金額</th></tr>
            </thead>
            <tbody>
              <tr v-for="emp in store.employeeReport!.employees" :key="emp.userId">
                <td>{{ emp.userName }}</td>
                <td class="r">{{ emp.count }}件</td>
                <td class="r">{{ formatAmount(emp.totalAmount) }}</td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="total-row">
                <td><strong>合計</strong></td>
                <td class="r"><strong>{{ store.employeeReport!.summary.totalCount }}件</strong></td>
                <td class="r"><strong>{{ formatAmount(store.employeeReport!.summary.grandTotal) }}</strong></td>
              </tr>
            </tfoot>
          </table>
        </section>

        <!-- 勘定科目別 合計 -->
        <section class="report-section">
          <h3>勘定科目別 合計</h3>
          <div v-if="!store.categoryReport?.categories?.length" class="empty-section">
            該当するデータはありません
          </div>
          <table v-else class="tbl" id="report-by-category">
            <thead>
              <tr><th>勘定科目</th><th class="r">件数</th><th class="r">合計金額</th><th class="r">構成比</th></tr>
            </thead>
            <tbody>
              <tr v-for="cat in store.categoryReport!.categories" :key="cat.categoryId">
                <td>{{ cat.categoryName }}</td>
                <td class="r">{{ cat.count }}件</td>
                <td class="r">{{ formatAmount(cat.totalAmount) }}</td>
                <td class="r">{{ formatPct(cat.percentage) }}</td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="total-row">
                <td><strong>合計</strong></td>
                <td class="r"><strong>{{ store.categoryReport!.summary.totalCount }}件</strong></td>
                <td class="r"><strong>{{ formatAmount(store.categoryReport!.summary.grandTotal) }}</strong></td>
                <td class="r"><strong>100.0%</strong></td>
              </tr>
            </tfoot>
          </table>
        </section>

        <!-- 集計対象の注記 -->
        <p class="report-note">※ 集計対象は「承認完了」ステータスのデータのみです</p>
      </template>
    </div>
  </AppLayout>
</template>

<style scoped>
.report-view__header{display:flex;justify-content:space-between;align-items:center;margin-bottom:var(--spacing-lg);flex-wrap:wrap;gap:var(--spacing-md)}
.report-view__header h2{font-size:1.25rem}
.year-month-selector{display:flex;align-items:center;gap:var(--spacing-sm)}
.year-month-selector label{font-size:.875rem;font-weight:500;color:var(--color-text-heading)}
.year-month-selector select{border:1px solid var(--color-input-border);border-radius:var(--radius-input);padding:var(--spacing-xs) var(--spacing-sm);font-size:.875rem}
.report-section{background:var(--color-surface-card);border-radius:var(--radius-card);box-shadow:var(--shadow-card);padding:var(--spacing-xl);margin-bottom:var(--spacing-lg)}
.report-section h3{font-size:1rem;margin-bottom:var(--spacing-md);padding-bottom:var(--spacing-sm);border-bottom:2px solid var(--color-primary)}
.tbl{width:100%;border-collapse:collapse}
.tbl th{background:var(--color-bg-light);padding:var(--spacing-sm) var(--spacing-md);text-align:left;font-size:.8125rem;font-weight:600;border-bottom:2px solid var(--color-border)}
.tbl td{padding:var(--spacing-sm) var(--spacing-md);border-bottom:1px solid var(--color-border);font-size:.875rem}
.tbl tbody tr:hover{background:var(--color-bg-light)}
.r{text-align:right}
.total-row{background:var(--color-bg-light);border-top:2px solid var(--color-border)}
.total-row td{padding:var(--spacing-sm) var(--spacing-md);font-size:.875rem}
.empty-section{text-align:center;padding:var(--spacing-lg);color:var(--color-text-muted);font-size:.875rem}
.report-note{font-size:.75rem;color:var(--color-text-muted);text-align:right;margin-top:var(--spacing-sm)}
.err{background:#fef2f2;border:1px solid var(--color-danger);color:var(--color-danger);padding:var(--spacing-sm) var(--spacing-md);border-radius:var(--radius-input);margin-bottom:var(--spacing-md)}
.empty{text-align:center;padding:var(--spacing-xxl);color:var(--color-text-muted)}
</style>
