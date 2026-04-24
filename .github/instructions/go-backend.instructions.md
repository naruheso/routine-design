---
applyTo: "backend/**/*.go"
description: "Go backend coding standards for 経費精算システム API. Applied to all Go files under backend/."
---

# Go Backend コーディング規約

## 基本方針

1. **型安全性**: sqlc による自動生成コード + 構造体タグバリデーション
2. **一貫性**: gofmt + golangci-lint で統一
3. **可読性**: 明確な命名と godoc コメント
4. **保守性**: 3層アーキテクチャに準拠
5. **テスト**: テーブル駆動テスト + testify

---

## アーキテクチャ

3層構造: **Handler → Service → Repository**

### Handler層 (`internal/handler/`)

**役割**: HTTPリクエスト/レスポンスの変換のみ

```go
// ✅ 良い例: バリデーション → サービス → レスポンス
func (h *ExpenseHandler) GetPaginated(c *gin.Context) {
    // 1. バリデーション
    params, err := validatePaginationParams(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
        return
    }

    // 2. サービス呼び出し
    result, err := h.service.FindPaginated(c.Request.Context(), params)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 3. レスポンス整形
    c.JSON(http.StatusOK, model.SuccessResponseWithPagination(result.Data, result.Pagination))
}
```

```go
// ❌ 悪い例: ビジネスロジックやDBアクセスを含む
func (h *ExpenseHandler) GetPaginated(c *gin.Context) {
    rows, err := h.db.Query("SELECT * FROM expenses") // NG: Handler でDB直接アクセス
    filtered := filterByCondition(rows)                // NG: Handler でビジネスロジック
}
```

**ルール**:
- ビジネスロジック禁止
- DBクエリ禁止
- Service層への委譲のみ
- エラーは共通の `handleError` メソッドで変換

### Service層 (`internal/service/`)

**役割**: ビジネスロジックの実装

```go
// ✅ 良い例: ビジネスロジックの集約
type ExpenseService struct {
    repo       ExpenseRepository
    categories CategoryRepository
}

// Submit は経費を申請する（下書き→承認待ちへ状態遷移）
func (s *ExpenseService) Submit(ctx context.Context, expenseID string, userID string) error {
    expense, err := s.repo.FindByID(ctx, expenseID)
    if err != nil {
        return fmt.Errorf("find expense %s: %w", expenseID, err)
    }
    // ビジネスルール: 申請者本人のみ申請可能
    if expense.UserID != userID {
        return fmt.Errorf("expense %s: %w", expenseID, ErrForbidden)
    }
    // ビジネスルール: 下書きのみ申請可能
    if expense.Status != "draft" {
        return fmt.Errorf("expense %s is not draft: %w", expenseID, ErrValidation)
    }
    return s.repo.UpdateStatus(ctx, expenseID, "pending_manager")
}
```

**ルール**:
- 単一責任（1サービス = 1ドメイン）
- 公開メソッドに godoc コメント必須
- 非公開メソッドで複雑ロジックを分割
- Repository はインターフェース経由で呼び出す

### Repository層 (`internal/repository/`)

**役割**: DB操作のラッパー

```go
// ✅ 良い例: sqlc 生成コードをラップ
type expenseRepository struct {
    q *sqlc.Queries
}

func (r *expenseRepository) FindByID(ctx context.Context, id string) (*model.Expense, error) {
    row, err := r.q.GetExpense(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("get expense: %w", err)
    }
    return toExpenseModel(row), nil
}
```

**ルール**:
- sqlc 生成コードのラップのみ
- `pgx.ErrNoRows` → ドメインエラーに変換
- 全てのエラーを `fmt.Errorf` でラップして返す

---

## バリデーション

### Handler層でのバリデーション

```go
// ✅ 良い例: 明確なバリデーション関数
func validatePaginationParams(c *gin.Context) (*PaginationParams, error) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    if page < 1 {
        page = 1
    }

    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    if limit < 1 || limit > 100 {
        limit = 20
    }

    sortBy := c.DefaultQuery("sortBy", "createdAt")
    allowedSort := map[string]bool{
        "createdAt": true, "expenseDate": true, "amount": true,
    }
    if !allowedSort[sortBy] {
        sortBy = "createdAt"
    }

    sortOrder := c.DefaultQuery("sortOrder", "desc")
    if sortOrder != "asc" {
        sortOrder = "desc"
    }

    return &PaginationParams{
        Page: page, Limit: limit,
        SortBy: sortBy, SortOrder: sortOrder,
        Search: c.Query("search"),
    }, nil
}
```

### 構造体バリデーション（POST/PUT）

Gin の `binding` タグを使用する:

```go
// ✅ リクエストボディのバリデーション
type CreateExpenseRequest struct {
    ExpenseDate string `json:"expenseDate" binding:"required"`
    CategoryID  string `json:"categoryId" binding:"required,uuid"`
    Amount      int    `json:"amount" binding:"required,min=1,max=999999"`
    Description string `json:"description" binding:"required,max=500"`
}

func (h *ExpenseHandler) Create(c *gin.Context) {
    var req CreateExpenseRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
        return
    }
    // ...
}
```

---

## 命名規則

### 変数・関数・型

| 種類 | 規則 | 例 |
|------|------|-----|
| パッケージ | 小文字単数形 | `handler`, `service`, `repository` |
| ファイル | スネークケース | `expense_service.go`, `expense_repo.go` |
| 公開型 | パスカルケース | `ExpenseService`, `PaginationParams` |
| 非公開型 | キャメルケース | `expenseRepository` |
| インターフェース | 用途名 or 動詞+er | `ExpenseRepository`, `Notifier` |
| 公開関数 | パスカルケース(動詞始まり) | `FindByID`, `CreateExpense` |
| 非公開関数 | キャメルケース(動詞始まり) | `findWithJoin`, `buildQuery` |
| 定数 | パスカルケース | `MaxPageSize`, `DefaultLimit` |
| 環境変数定数 | パスカルケース | `EnvDBHost`, `EnvPort` |
| テストファイル | `*_test.go` | `expense_service_test.go` |

### データベース

| 種類 | 規則 | 例 |
|------|------|-----|
| テーブル名 | snake_case（複数形） | `expenses`, `approval_histories` |
| カラム名 | snake_case | `expense_date`, `created_at` |
| インデックス名 | `idx_table_column` | `idx_expenses_user_id` |

### JSON タグ

```go
// ✅ キャメルケースで出力（フロントエンド互換）
type Expense struct {
    ID          string    `json:"id"`
    UserID      string    `json:"userId"`
    ExpenseDate string    `json:"expenseDate"`
    Amount      int       `json:"amount"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

---

## エラーハンドリング

### センチネルエラー定義

```go
// internal/service/errors.go
var (
    ErrNotFound   = errors.New("not found")
    ErrValidation = errors.New("validation error")
    ErrConflict   = errors.New("conflict")
    ErrForbidden  = errors.New("forbidden")
)
```

### Service層: エラーのラップ

```go
// ✅ 良い例: コンテキスト情報を付与してラップ
func (s *ExpenseService) FindByID(ctx context.Context, id string) (*Expense, error) {
    expense, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("find expense %s: %w", id, err)
    }
    return expense, nil
}
```

### Handler層: エラー → HTTPステータス変換

```go
// ✅ 良い例: 共通エラーハンドラ
func (h *BaseHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrNotFound):
        c.JSON(http.StatusNotFound, model.ErrorResponse("not found"))
    case errors.Is(err, service.ErrValidation):
        c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
    case errors.Is(err, service.ErrConflict):
        c.JSON(http.StatusConflict, model.ErrorResponse("already exists"))
    case errors.Is(err, service.ErrForbidden):
        c.JSON(http.StatusForbidden, model.ErrorResponse("forbidden"))
    default:
        slog.Error("internal error", "error", err)
        c.JSON(http.StatusInternalServerError, model.ErrorResponse("internal server error"))
    }
}
```

```go
// ❌ 悪い例: エラー詳細をそのまま返す
c.JSON(500, gin.H{"error": err.Error()}) // 内部情報漏洩リスク
```

---

## レスポンス形式

`model/common.go` の共通ヘルパーを必ず使う:

```go
// 成功レスポンス
model.SuccessResponse(data)                           // {"success":true,"data":...}
model.SuccessResponseWithCount(data, count)           // {"success":true,"data":...,"count":N}
model.SuccessResponseWithPagination(data, pagination) // {"success":true,"data":...,"pagination":{...}}

// エラーレスポンス
model.ErrorResponse(message)                          // {"success":false,"error":"..."}
```

### ページネーションレスポンス

```go
// フロントエンド互換のJSON形式
type Pagination struct {
    CurrentPage int  `json:"currentPage"`
    TotalPages  int  `json:"totalPages"`
    TotalCount  int  `json:"totalCount"`
    Limit       int  `json:"limit"`
    HasNext     bool `json:"hasNext"`
    HasPrev     bool `json:"hasPrev"`
}
```

---

## 依存性注入

- Service はコンストラクタでインターフェースを受け取る
- Handler はコンストラクタで Service を受け取る
- `main.go` で全ての依存関係を組み立てる（DIコンテナ不使用）

```go
// cmd/api/main.go
repo := repository.NewExpenseRepository(db)
svc := service.NewExpenseService(repo)
h := handler.NewExpenseHandler(svc)
```

### インターフェース定義

```go
// ✅ インターフェースは利用側（Service）で定義
// internal/service/expense_service.go
type ExpenseRepository interface {
    FindAll(ctx context.Context, limit int) ([]model.Expense, error)
    FindByID(ctx context.Context, id string) (*model.Expense, error)
    FindPaginated(ctx context.Context, params PaginationParams) ([]model.Expense, int, error)
    Create(ctx context.Context, expense model.Expense) error
    UpdateStatus(ctx context.Context, id string, status string) error
}
```

---

## ルーティング

Gin のルートグループを使用して整理する:

```go
// cmd/api/main.go
func setupRouter(h *Handlers) *gin.Engine {
    r := gin.Default()

    // CORS ミドルウェア
    r.Use(corsMiddleware())

    // 公開エンドポイント
    r.GET("/health", h.health.Check)

    api := r.Group("/api")
    {
        // 認証
        api.POST("/auth/login", h.auth.Login)

        // 認証必須エンドポイント
        authorized := api.Group("/")
        authorized.Use(authMiddleware())
        {
            // 経費
            expenses := authorized.Group("/expenses")
            {
                expenses.GET("", h.expense.GetPaginated)
                expenses.POST("", h.expense.Create)
                expenses.GET("/:id", h.expense.GetByID)
                expenses.PUT("/:id", h.expense.Update)
                expenses.POST("/:id/submit", h.expense.Submit)
            }

            // 承認
            approvals := authorized.Group("/approvals")
            {
                approvals.GET("/pending", h.approval.GetPending)
                approvals.POST("/:id/approve", h.approval.Approve)
                approvals.POST("/:id/return", h.approval.Return)
                approvals.POST("/:id/reject", h.approval.Reject)
            }
        }
    }

    return r
}
```

---

## ミドルウェア

### JWT 認証ミドルウェア

```go
// internal/middleware/auth.go
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authorization ヘッダーからトークン取得
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse("認証が必要です"))
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateToken(tokenString, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse("無効なトークンです"))
            return
        }

        // コンテキストにユーザー情報をセット
        c.Set("userID", claims.UserID)
        c.Set("roles", claims.Roles)
        c.Next()
    }
}
```

### CORS ミドルウェア

```go
// internal/middleware/cors.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        allowedOrigins := []string{"http://localhost:5173"}

        for _, allowed := range allowedOrigins {
            if origin == allowed {
                c.Header("Access-Control-Allow-Origin", origin)
                break
            }
        }
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}
```

---

## トランザクション

```go
// ✅ 複数操作は必ずトランザクション
func (r *expenseRepository) CreateWithHistory(ctx context.Context, expense model.Expense) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback(ctx)

    qtx := r.q.WithTx(tx)

    if err := qtx.InsertExpense(ctx, toInsertParams(expense)); err != nil {
        return fmt.Errorf("insert expense: %w", err)
    }

    return tx.Commit(ctx)
}
```

---

## ログ

- `log/slog` (標準ライブラリ) を使用
- リクエストログはミドルウェアで自動出力
- Service/Repository 層では構造化ログを使う

```go
// ✅ 良い例: 構造化ログ
slog.Info("経費申請完了", "expenseId", id, "userId", userID, "status", "pending_manager")
slog.Error("経費取得失敗", "expenseId", id, "error", err)
```

```go
// ❌ 悪い例
log.Println("error:", err)                           // 非構造化
slog.Info("JWT Secret: " + secret)                   // 機密情報漏洩
fmt.Printf("debug: %v\n", data)                     // デバッグプリント残留
```

---

## SQL / sqlc

- 全てのクエリは `sql/queries/*.sql` に定義
- `sqlc generate` で Go コードを自動生成
- 手書きSQLは Repository 層に書かない（sqlc 経由のみ）
- パラメータ化クエリを必ず使う（SQLインジェクション防止）

```sql
-- sql/queries/expenses.sql

-- name: GetExpense :one
SELECT * FROM expenses WHERE id = $1;

-- name: ListExpensesByUser :many
SELECT * FROM expenses WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2;

-- name: InsertExpense :exec
INSERT INTO expenses (user_id, expense_date, category_id, amount, description, status)
VALUES ($1, $2, $3, $4, $5, $6);
```

---

## コンテキスト

- 全ての DB 操作・外部API呼び出しに `context.Context` を渡す
- Handler は `c.Request.Context()` から取得
- タイムアウトが必要な場合は `context.WithTimeout` を使用

```go
// ✅ 良い例
func (h *ExpenseHandler) Create(c *gin.Context) {
    ctx := c.Request.Context()
    result, err := h.service.Create(ctx, req)
    // ...
}
```

---

## import 順序

goimports に従う。3グループに分けて空行で区切る:

```go
import (
    // 1. 標準ライブラリ
    "context"
    "fmt"
    "net/http"
    "time"

    // 2. 外部ライブラリ
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"

    // 3. 内部パッケージ
    "github.com/amical/routine-design/backend/internal/model"
    "github.com/amical/routine-design/backend/internal/service"
)
```

---

## コメント規約

### godoc コメント

```go
// ExpenseService は経費ドメインのビジネスロジックを提供する。
type ExpenseService struct { ... }

// FindPaginated はページネーション付きの経費一覧を返す。
// ステータスフィルタが指定された場合、該当ステータスの経費のみを返す。
func (s *ExpenseService) FindPaginated(ctx context.Context, params PaginationParams) (*PaginatedResult, error) {
```

### インラインコメント

```go
// ✅ 良い例: 「なぜ」を説明
// 承認フローの整合性のため、ステータス更新と履歴挿入をトランザクションで実行
if err := s.repo.ApproveWithHistory(ctx, id, actorID); err != nil {
    return err
}

// ❌ 悪い例: コードと同じことを書く
// IDを取得する
id := c.Param("id")
```

---

## セキュリティ

- 入力値は Handler 層で必ずバリデーション（`binding` タグ + カスタムバリデーション）
- SQL は sqlc のパラメータ化クエリのみ（SQLインジェクション防止）
- エラーレスポンスに内部情報を含めない
- APIキー・パスワードをログやレスポンスに出力しない
- CORS は許可オリジンを明示的に設定（ワイルドカード禁止）
- JWT シークレットは環境変数で管理

---

## NULL ハンドリング

NULLableカラムには pgx の Nullable 型を使う:

```go
// ✅ NULL可能フィールド
type Expense struct {
    ID               string          `json:"id"`
    Description      string          `json:"description"`
    ReceiptImagePath pgtype.Text     `json:"receiptImagePath"` // NULL可能
    SubmittedAt      pgtype.Timestamp `json:"submittedAt"`     // NULL可能
}

// JSON出力時: null → JSON null, Valid → 値
```

---

## 日時フォーマット

- DB ↔ Go: `time.Time` (pgx が自動変換)
- Go → JSON: `time.RFC3339` (`2026-04-15T12:00:00Z`) で出力
- フロントエンドとの互換性のため ISO8601 統一
