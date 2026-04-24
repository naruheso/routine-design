<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ROLE_DISPLAY_NAMES } from '@/types/user'

const authStore = useAuthStore()
const router = useRouter()

// 集計レポートは expense_admin / finance_director のみ表示
const isReportViewer = computed(() => {
  return authStore.hasRole('expense_admin') || authStore.hasRole('finance_director')
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="app-header__left">
        <h1 class="app-header__title">経費精算システム</h1>
      </div>
      <nav class="app-header__nav">
        <router-link to="/" class="nav-link">経費一覧</router-link>
        <router-link v-if="authStore.isApprover" to="/approvals" class="nav-link">承認待ち</router-link>
        <router-link v-if="isReportViewer" to="/reports" class="nav-link">集計レポート</router-link>
      </nav>
      <div class="app-header__right">
        <span class="app-header__user">
          {{ authStore.userName }}
          <span class="app-header__role">
            ({{ authStore.userRoles.map(r => ROLE_DISPLAY_NAMES[r] ?? r).join(', ') }})
          </span>
        </span>
        <button class="btn-logout" @click="handleLogout">ログアウト</button>
      </div>
    </header>
    <main class="app-main">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: var(--color-bg-base);
  border-bottom: 1px solid var(--color-border);
  padding: 0 var(--spacing-xl);
  height: 56px;
  display: flex;
  align-items: center;
  gap: var(--spacing-xl);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.app-header__title {
  font-size: 1.125rem;
  color: var(--color-heading-blue);
  white-space: nowrap;
}

.app-header__nav {
  display: flex;
  gap: var(--spacing-md);
  flex: 1;
}

.nav-link {
  color: var(--color-text-body);
  font-weight: 500;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-button);
  transition: all var(--transition-fast);
}

.nav-link:hover {
  background: var(--color-bg-light);
  color: var(--color-primary);
}

.nav-link.router-link-active {
  color: var(--color-primary);
  background: rgba(40, 100, 240, 0.08);
}

.app-header__right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.app-header__user {
  font-size: 0.875rem;
  color: var(--color-text-body);
}

.app-header__role {
  color: var(--color-text-muted);
  font-size: 0.75rem;
}

.btn-logout {
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-button);
  padding: var(--spacing-xs) var(--spacing-md);
  font-size: 0.8125rem;
  color: var(--color-text-body);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-logout:hover {
  background: var(--color-bg-light);
  border-color: var(--color-text-muted);
}

.app-main {
  flex: 1;
  padding: var(--spacing-xl);
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}
</style>
