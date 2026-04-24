<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const showError = ref(false)

async function handleLogin() {
  showError.value = false
  const success = await authStore.login(email.value, password.value)
  if (success) {
    router.push('/')
  } else {
    showError.value = true
  }
}
</script>

<template>
  <div class="login-view">
    <div class="login-card">
      <div class="login-card__header">
        <h1>経費精算システム</h1>
        <p>Routine Design</p>
      </div>

      <form class="login-card__form" @submit.prevent="handleLogin">
        <div v-if="showError && authStore.error" class="login-card__error">
          {{ authStore.error }}
        </div>

        <div class="form-group">
          <label for="email">メールアドレス</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="例: taro.ippan@example.com"
            required
            autocomplete="email"
          />
        </div>

        <div class="form-group">
          <label for="password">パスワード</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="パスワードを入力"
            required
            autocomplete="current-password"
          />
        </div>

        <button
          type="submit"
          class="btn-primary"
          :disabled="authStore.loading"
        >
          {{ authStore.loading ? 'ログイン中...' : 'ログイン' }}
        </button>
      </form>

      <div class="login-card__hint">
        <p><strong>テストアカウント:</strong></p>
        <ul>
          <li>申請者: taro.ippan@example.com</li>
          <li>上長: ichiro.manager@example.com</li>
          <li>経費担当: tanto.keiri@example.com</li>
          <li>経理部長: bucho.keiri@example.com</li>
        </ul>
        <p>パスワード: password123</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, var(--color-primary-darkest) 0%, var(--color-primary) 100%);
}

.login-card {
  background: var(--color-surface-card);
  border-radius: var(--radius-dialog);
  box-shadow: var(--shadow-floating);
  padding: var(--spacing-xxl);
  width: 100%;
  max-width: 420px;
  margin: var(--spacing-md);
}

.login-card__header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.login-card__header h1 {
  font-size: 1.5rem;
  color: var(--color-heading-blue);
  margin-bottom: var(--spacing-xs);
}

.login-card__header p {
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.login-card__error {
  background: #fef2f2;
  border: 1px solid var(--color-danger);
  color: var(--color-danger);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-input);
  font-size: 0.875rem;
  margin-bottom: var(--spacing-md);
}

.login-card__form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-heading);
}

.form-group input {
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-input);
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: 1rem;
  transition: border-color var(--transition-fast);
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(40, 100, 240, 0.15);
}

.btn-primary {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border: 2px solid var(--color-primary);
  border-radius: var(--radius-button);
  padding: var(--spacing-sm) var(--spacing-lg);
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-card__hint {
  margin-top: var(--spacing-xl);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--color-border);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.login-card__hint ul {
  margin: var(--spacing-xs) 0;
  padding-left: var(--spacing-md);
}

.login-card__hint li {
  margin-bottom: 2px;
}
</style>
