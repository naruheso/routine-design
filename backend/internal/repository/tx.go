package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const txKey = contextKey("tx")

// DBQuerier は pgxpool.Pool と pgx.Tx の共通インターフェース。
type DBQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TransactionManager はトランザクションの境界を管理する。
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

// NewTransactionManager はTransactionManagerを生成する。
func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &txManager{pool: pool}
}

// RunInTransaction は渡された関数をトランザクション内で実行する。
// 既に関数がトランザクションコンテキストで呼ばれた場合は、そのまま実行する。
func (tm *txManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 既にトランザクション内ならそのまま実行
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	// パニックまたはエラー時にも確実にロールバック (Commit済みの場合は何もしない)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetDBQuerier はcontextからトランザクション(Tx)を取得するか、通常のPoolを返す。
// Repository内の各メソッドはこれを使ってクエリを実行する。
func GetDBQuerier(ctx context.Context, pool *pgxpool.Pool) DBQuerier {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return pool
}
