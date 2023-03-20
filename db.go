package go_orm

import (
	"database/sql"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	r  *registry
	db *sql.DB
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
		r: &registry{
			models: map[reflect.Type]*model{},
		},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r *registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}
