package go_orm

import "reflect"

type DBOption func(db *DB)

type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{
			models: map[reflect.Type]*model{},
		},
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
