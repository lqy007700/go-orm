package valuer

import (
	"database/sql"
	err2 "go-orm/internal/err"
	go_orm "go-orm/internal/model"
	"reflect"
)

type ReflectValue struct {
	t     any
	model *go_orm.Model
}

func NewReflectValue(t any, model *go_orm.Model) *ReflectValue {
	return &ReflectValue{
		t:     t,
		model: model,
	}
}

func (r *ReflectValue) SetColumns(rows *sql.Rows) error {
	if !rows.Next() {
		return err2.ErrNoRows
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(columns) > len(r.model.FieldMap) {
		return err2.ErrTooManyReturnedColumns
	}

	vals := make([]any, 0, len(columns))
	eleVals := make([]reflect.Value, 0, len(columns))
	for _, column := range columns {
		f := r.model.ColumnMap[column]

		fdVal := reflect.New(f.Typ)
		eleVals = append(eleVals, fdVal.Elem())

		vals = append(vals, fdVal.Interface())
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	t := r.t
	tVal := reflect.ValueOf(t).Elem()
	for i, column := range columns {
		f, ok := r.model.ColumnMap[column]
		if !ok {
			return err2.NewErrUnknownColumn(column)
		}
		tVal.FieldByName(f.GoName).Set(eleVals[i])
	}
	return nil
}
