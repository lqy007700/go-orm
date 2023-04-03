package go_orm

import (
	"context"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-orm/internal/model"
	"go-orm/internal/valuer"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	db *sql.DB
	core
}

func Open(driver, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		db: db,

		core: core{
			r: &model.Registrys{
				Models: map[reflect.Type]*model.Model{},
			},
			valCreator: valuer.NewUnsafeValue,
			dialect:    &mysqlDialect{},
		},
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (db *DB) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args)
}

func (db *DB) exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args)
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

func DBWithRegistry(r *model.Registrys) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBWithDialect(d Dialect) DBOption {
	return func(db *DB) {
		db.dialect = d
	}
}

func DBWithMiddleware(ms ...MiddleWare) DBOption {
	return func(db *DB) {
		db.ms = ms
	}
}
