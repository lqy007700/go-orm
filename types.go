package go_orm

import (
	"context"
	"database/sql"
)

// QueryI select语句
type QueryI[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

type Executor[T any] interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	args []any
}
