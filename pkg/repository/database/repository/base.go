package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Base struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Base {
	return &Base{pool: pool}
}

func (b *Base) Begin(ctx context.Context, f func(pgx.Tx) error) error {
	if err := b.pool.BeginFunc(ctx, f); err != nil {
		return fmt.Errorf("couldn't wrap into txn: %w", err)
	}

	return nil
}

func (b *Base) BeginTx(ctx context.Context, txOptions pgx.TxOptions, f func(tx pgx.Tx) error) error {
	if err := b.pool.BeginTxFunc(ctx, txOptions, f); err != nil {
		return fmt.Errorf("couldn't wrap into txn: %w", err)
	}

	return nil
}
