---
applyTo: "frontend/**/*.{vue,ts}"
description: "Vue 3 frontend coding standards for pachi-admin. Applied to all Vue/TypeScript files under frontend/."
---

# Vue Frontend コーディング規約

## 基本方針

1. **Composition API**: `<script setup lang="ts">` を標準とする
2. **型安全性**: TypeScript strict mode + 明示的な型定義
3. **一貫性**: ESLint + Prettier で統一
4. **可読性**: 単一責務のコンポーネント + 明確な命名
5. **テスト**: Vitest + @vue/test-utils

---

## アーキテクチャ

4層構造: **View → Component → Store (Pinia) → API Client**

### View層 (`src/views/`)

**役割**: ページ単位のルートコンポーネント。データ取得の起点。

```vue
<!-- ✅ 良い例: Store経由でデータ取得、表示はComponentに委譲 -->
<script setup lang="ts">
import { onMounted } from 'vue'
import { useChannelStore } from '@/stores/channel'
import ChannelTable from '@/components/channel/ChannelTable.vue'
import Pagination from '@/components/common/Pagination.vue'

const store = useChannelStore()

onMounted(() => {
  store.loadChannels()
})
</script>

<template>
  <div class="channel-list-view">
    <h1>チャンネル管理</h1>
    <ChannelTable :channels="store.channels" :loading="store.loading" />
    <Pagination
      :current-page="store.currentPage"
      :total-pages="store.totalPages"
      @page-change="store.changePage"
    />
  </div>
</template>
```

```vue
<!-- ❌ 悪い例: Viewで直接fetch、ロジックを含む -->
<script setup lang="ts">
const response = await fetch('/api/channels') // NG: View で直接API呼び出し
const filtered = channels.filter(...)          // NG: View でビジネスロジック
</script>
```

**ルール**:
- 直接の `fetch` 呼び出し禁止（Store / API Client 経由）
- ビジネスロジック禁止（Store に委譲）
- レイアウトと子コンポーネント配置のみ

### Component層 (`src/components/`)

**役割**: 再利用可能なUI部品。Props と Emits で通信。

```vue
<!-- ✅ 良い例: Props で受け取り、Emits でイベント通知 -->
<script setup lang="ts">
interface Props {
  channels: Channel[]
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  select: [channelId: string]
  delete: [channelId: string]
}>()
</script>
```

**ルール**:
- Props はインターフェースで型定義（`defineProps<Props>()`）
- Emits は型付き定義（`defineEmits<{...}>()`）
- Store への直接依存は避ける（Props/Emits で通信）
- `common/` 配下はドメイン非依存の汎用コンポーネント

### Store層 (`src/stores/`)

**役割**: 状態管理 + API 呼び出しの橋渡し

```typescript
// ✅ 良い例: Pinia Setup Store
import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { fetchChannelsPaginated } from '@/api/channels'
import type { Channel, PaginationParams } from '@/types/channel'

export const useChannelStore = defineStore('channel', () => {
  // State
  const channels = ref<Channel[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const currentPage = ref(1)
  const totalPages = ref(1)

  // Getters
  const hasChannels = computed(() => channels.value.length > 0)

  // Actions
  async function loadChannels(params?: Partial<PaginationParams>) {
    loading.value = true
    error.value = null
    try {
      const result = await fetchChannelsPaginated({
        page: currentPage.value,
        limit: 50,
        ...params,
      })
      if (result.success) {
        channels.value = result.data
        totalPages.value = result.pagination?.totalPages ?? 1
      } else {
        error.value = result.error ?? 'データ取得に失敗しました'
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '不明なエラー'
    } finally {
      loading.value = false
    }
  }

  function changePage(page: number) {
    currentPage.value = page
    loadChannels()
  }

  return { channels, loading, error, currentPage, totalPages, hasChannels, loadChannels, changePage }
})
```

**ルール**:
- Setup Store 形式（`defineStore('name', () => {...})`）を使用
- State は `ref()` / `reactive()` で定義
- Getters は `computed()` で定義
- Actions は通常の関数として定義
- API 呼び出しは API Client 経由

### API Client層 (`src/api/`)

**役割**: Go backend API との通信

```typescript
// ✅ 良い例: 型安全な API クライアント
import { apiClient } from './client'
import type { Channel, PaginationParams } from '@/types/channel'
import type { ApiResponse, PaginatedResponse } from '@/types/api'

export async function fetchChannelsPaginated(
  params: PaginationParams
): Promise<PaginatedResponse<Channel[]>> {
  const query = new URLSearchParams({
    page: params.page.toString(),
    limit: params.limit.toString(),
    sortBy: params.sortBy ?? 'createdAt',
    sortOrder: params.sortOrder ?? 'desc',
  })
  if (params.search) query.set('search', params.search)

  return apiClient.get<Channel[]>(`/channels/paginated?${query}`)
}

export async function deleteChannel(channelId: string): Promise<ApiResponse<void>> {
  return apiClient.delete(`/channels/${channelId}`)
}
```

**ルール**:
- 全関数に戻り値型を明示
- クエリパラメータは `URLSearchParams` で構築
- パスパラメータはテンプレートリテラルで埋め込み

---

## 命名規則

### ファイル・ディレクトリ

| 種類 | 規則 | 例 |
|------|------|-----|
| View | パスカルケース + View | `ChannelListView.vue` |
| Component | パスカルケース | `ChannelTable.vue`, `Pagination.vue` |
| Store | キャメルケース | `channel.ts` → `useChannelStore` |
| API | キャメルケース | `channels.ts` |
| 型定義 | キャメルケース | `channel.ts` |
| Composable | キャメルケース + use | `usePagination.ts` |
| テスト | 元ファイル名 + .test | `ChannelTable.test.ts` |
| CSS | scoped style（SFC内） | `<style scoped>` |

### コード

| 種類 | 規則 | 例 |
|------|------|-----|
| コンポーネント名 | パスカルケース | `ChannelTable`, `LoadingSpinner` |
| Props | キャメルケース | `channelId`, `isLoading` |
| Emits | キャメルケース | `pageChange`, `itemSelect` |
| Store | use + パスカルケース + Store | `useChannelStore` |
| Composable | use + パスカルケース | `usePagination`, `useDebounce` |
| 型/インターフェース | パスカルケース | `Channel`, `PaginationParams` |
| 定数 | UPPER_SNAKE_CASE | `API_BASE_URL`, `MAX_PAGE_SIZE` |
| CSS クラス | ケバブケース | `channel-list`, `sort-button` |

### テンプレート内

```vue
<!-- ✅ 良い例 -->
<ChannelTable :channel-id="channelId" @page-change="handlePageChange" />

<!-- ❌ 悪い例 -->
<channelTable :channelId="channelId" @pageChange="handlePageChange" />
```

- コンポーネント: パスカルケース（`<ChannelTable />`）
- Props/Events: ケバブケース（`:channel-id`, `@page-change`）

---

## 型定義

### API レスポンス共通型

```typescript
// src/types/api.ts
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
```

### ドメイン型

```typescript
// src/types/channel.ts
export interface Channel {
  id: number
  channelId: string
  title: string
  link: string | null
  authorName: string | null
  published: string | null
  createdAt: string | null
  counts: ChannelCounts | null
  labels: Label[]
}

export interface PaginationParams {
  page: number
  limit: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  search?: string
}
```

**ルール**:
- `any` 使用禁止（`unknown` を使う）
- API レスポンスの null 可能フィールドは `T | null` で明示
- 日時は `string`（ISO8601）で受け取り、表示時にフォーマット

---

## Composables

### 基本パターン

```typescript
// ✅ 良い例: 引数と戻り値の型を明示
import { ref, watch, type Ref } from 'vue'

export function useDebounce<T>(source: Ref<T>, delay: number = 300): Ref<T> {
  const debounced = ref(source.value) as Ref<T>
  let timer: ReturnType<typeof setTimeout>

  watch(source, (val) => {
    clearTimeout(timer)
    timer = setTimeout(() => {
      debounced.value = val
    }, delay)
  })

  return debounced
}
```

**ルール**:
- `use` プレフィックス必須
- リアクティブな値を返す
- クリーンアップが必要な場合は `onUnmounted` で解除

---

## エラーハンドリング

### API クライアント共通

```typescript
// src/api/client.ts
const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api'

async function request<T>(path: string, options?: RequestInit): Promise<ApiResponse<T>> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    })
    if (!res.ok) {
      const body = await res.json().catch(() => null)
      return { success: false, error: body?.error ?? `HTTP ${res.status}` }
    }
    return await res.json()
  } catch (e) {
    return { success: false, error: e instanceof Error ? e.message : '通信エラー' }
  }
}
```

### コンポーネントでのエラー表示

```vue
<!-- ✅ 良い例: 統一的なエラー表示 -->
<template>
  <ErrorMessage v-if="store.error" :message="store.error" />
  <LoadingSpinner v-else-if="store.loading" />
  <div v-else>
    <!-- コンテンツ -->
  </div>
</template>
```

---

## レスポンス形式

Go backend API のレスポンス形式に準拠する:

```json
{
  "success": true,
  "data": [...],
  "count": 10,
  "pagination": {
    "currentPage": 1,
    "totalPages": 5,
    "totalCount": 100,
    "limit": 20,
    "hasNext": true,
    "hasPrev": false
  }
}
```

---

## ルーティング

### 遅延ロード

```typescript
// ✅ 良い例: ルート単位でコード分割
const routes = [
  {
    path: '/channels',
    name: 'channels',
    component: () => import('@/views/ChannelListView.vue'),
  },
]
```

### URL クエリパラメータ同期

```typescript
// ✅ 良い例: useUrlParams composable で同期
const route = useRoute()
const router = useRouter()

// URLから初期値を読み取り
const page = ref(Number(route.query.page) || 1)

// 状態変更時にURLを更新
watch(page, (val) => {
  router.replace({ query: { ...route.query, page: val.toString() } })
})
```

---

## デザインシステム

> **重要**: 本プロジェクトのUIデザインは [DESIGN.md](./DESIGN.md)（freee Vibes Design System 準拠）に従うこと。
> CSS 変数・色・フォント・スペーシング等は DESIGN.md のトークン値を正とする。

### デザイントークンの適用方針

- CSS 変数は `assets/styles/variables.css` で一元管理し、DESIGN.md のトークン値を反映する
- コンポーネント内でハードコードされた色・サイズを使用しない（必ず CSS 変数経由）
- 新しい色やサイズが必要な場合は、まず DESIGN.md のパレットから選択する

### CSS 変数と DESIGN.md トークンの対応

```css
:root {
  /* === Primary（DESIGN.md §2 Primary） === */
  --color-primary: #2864f0;
  --color-primary-hover: #285ac8;
  --color-primary-dark: #1e46aa;
  --color-primary-darkest: #143278;

  /* === Semantic（DESIGN.md §2 Semantic） === */
  --color-danger: #dc1e32;
  --color-danger-hover: #a51428;
  --color-warning: #ffb91e;
  --color-success: #00963c;

  /* === Text（DESIGN.md §2 Neutral） === */
  --color-text-heading: #323232;
  --color-text-body: #595959;
  --color-text-muted: #8c8989;
  --color-text-inverse: #ffffff;
  --color-heading-blue: #1e46aa;

  /* === Surface & Borders（DESIGN.md §2 Surface） === */
  --color-bg-base: #ffffff;
  --color-bg-light: #f7f5f5;
  --color-bg-lighter: #f0eded;
  --color-surface-card: #ffffff;
  --color-border: #e9e7e7;
  --color-input-border: #cccccc;

  /* === Spacing（DESIGN.md §5 Spacing Scale） === */
  --spacing-xs: 0.25rem;   /* 4px */
  --spacing-sm: 0.5rem;    /* 8px */
  --spacing-md: 1rem;      /* 16px */
  --spacing-lg: 1.5rem;    /* 24px */
  --spacing-xl: 2rem;      /* 32px */
  --spacing-xxl: 3rem;     /* 48px */

  /* === Border Radius（DESIGN.md §5 Radius Scale） === */
  --radius-input: 4px;
  --radius-button: 8px;
  --radius-card: 0.75rem;  /* 12px */
  --radius-floating: 1rem; /* 16px */
  --radius-dialog: 1.5rem; /* 24px */

  /* === Shadow（DESIGN.md §6 Depth） === */
  --shadow-card: 0 0 1rem rgba(0,0,0,0.1), 0 0.125rem 0.25rem rgba(0,0,0,0.2);
  --shadow-floating: 0 0 1.5rem rgba(0,0,0,0.1), 0 0.25rem 0.5rem rgba(0,0,0,0.2);

  /* === Transition === */
  --transition-fast: 0.2s;
  --transition-standard: 0.3s;
}
```

### フォント戦略

本プロジェクトは**プロダクトUI**として構築するため、DESIGN.md §3 の Vibes システムフォントスタックを使用する。

```css
body {
  font-family: '-apple-system', BlinkMacSystemFont, 'Helvetica Neue',
    'ヒラギノ角ゴ ProN', 'Hiragino Kaku Gothic ProN', Arial,
    'メイリオ', Meiryo, sans-serif;
  font-size: 0.875rem; /* 14px — Vibes Body 標準 */
  line-height: 1.5;
  color: var(--color-text-body);
}
```

> **禁止**: Noto Sans JP 等の Web フォントをプロダクト UI に使用しない（DESIGN.md §7 Don't）

### コンポーネントスタイリング（DESIGN.md §4 準拠）

**ボタン**:
```css
.btn-primary {
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border: 2px solid var(--color-primary);
  border-radius: var(--radius-button);
  font-weight: 500;
}
.btn-primary:hover { background: var(--color-primary-hover); }
```

**入力欄**:
```css
.input {
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius-input);
  font-size: 1rem;
}
```

**カード**:
```css
.card {
  background: var(--color-surface-card);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
}
```

### CSS ルール

### Scoped Style

```vue
<style scoped>
.channel-list {
  padding: var(--spacing-md);
}

.channel-list__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

**ルール**:
- `<style scoped>` を標準とする
- BEM ライクな命名（`block__element--modifier`）
- CSS 変数はグローバル (`assets/styles/variables.css`) で定義し、DESIGN.md のトークンに準拠
- コンポーネント内でハードコード値（`#2864f0` 等）を直接書かない。必ず CSS 変数を使用
- スペーシングは 4px の倍数に揃える（`--spacing-*` 変数を使用）
- テキスト色に `#000000` を使用しない（`--color-text-heading: #323232` を使用）
- `!important` 使用禁止

---

## セキュリティ

- `v-html` 使用禁止（XSS 防止）。必要な場合は DOMPurify でサニタイズ
- API キーをフロントエンドコードに含めない
- 環境変数は `VITE_` プレフィックス付きのみ使用
- ユーザー入力はバリデーションしてから API に送信

---

## import 順序

```typescript
// 1. Vue / ライブラリ
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

// 2. Store / API
import { useChannelStore } from '@/stores/channel'
import { fetchChannels } from '@/api/channels'

// 3. コンポーネント
import ChannelTable from '@/components/channel/ChannelTable.vue'
import Pagination from '@/components/common/Pagination.vue'

// 4. 型
import type { Channel, PaginationParams } from '@/types/channel'

// 5. ユーティリティ
import { formatDate, formatNumber } from '@/utils/format'
```

---

## SFC 構造順序

```vue
<script setup lang="ts">
// 1. imports
// 2. props / emits
// 3. store / composables
// 4. reactive state (ref, reactive)
// 5. computed
// 6. watch
// 7. lifecycle hooks (onMounted, etc.)
// 8. methods
</script>

<template>
  <!-- HTML -->
</template>

<style scoped>
/* CSS */
</style>
```
