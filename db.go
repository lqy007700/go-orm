package go_orm

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-orm/internal/model"
	"go-orm/internal/valuer"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	r  *model.Registrys
	db *sql.DB

	valCreator valuer.Creator
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
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
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
