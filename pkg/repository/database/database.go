package database

import (
	"context"
	"recipe-app/pkg/domain/constant"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v4"
)

type Beginner interface {
	Begin(ctx context.Context, f func(tx pgx.Tx) error) error
}

type BeginTxer interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions, f func(tx pgx.Tx) error)
}

type Transactable interface {
	Beginner
	BeginTxer
}

type Tabler interface {
	Table() constant.Table
}

func wrapSelectPagedCompose(page, size uint64, query *sq.SelectBuilder) sq.SelectBuilder {
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * size

	return query.Limit(size).Offset(offset)
}
