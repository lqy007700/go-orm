package go_orm

import (
	"context"
	"database/sql"
	"go-orm/internal/model"
	"go-orm/internal/valuer"
)

type Tx struct {
	tx *sql.Tx
	core
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) getCore() core {
	return t.core
}

func (t *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args)
}

func (t *Tx) exec(query string, args ...any) (sql.Result, error) {
	return t.tx.Exec(query, args)
}

type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	exec(query string, args ...any) (sql.Result, error)
}

type core struct {
	r          model.Registry
	valCreator valuer.Creator
	dialect    Dialect
	ms         []MiddleWare
}
