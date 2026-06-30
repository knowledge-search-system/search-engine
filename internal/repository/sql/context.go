package sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func ExecutorFromContext(ctx context.Context, pool Executor) Executor {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return pool
}
