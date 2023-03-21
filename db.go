package go_orm

import (
	"database/sql"
	"go-orm/internal/model"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	r  *model.Registrys
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
		r: &model.Registrys{
			Models: map[reflect.Type]*model.Model{},
		},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r *model.Registrys) DBOption {
	return func(db *DB) {
		db.r = r
	}
}
