package valuer

import (
	"database/sql"
	err2 "go-orm/internal/err"
	"go-orm/internal/model"
	"reflect"
	"unsafe"
)

type UnsafeValue struct {
	t     any
	model *model.Model
	addr  unsafe.Pointer
}

func NewUnsafeValue(t any, model *model.Model) Value {
	addr := unsafe.Pointer(reflect.ValueOf(t).Pointer())
	return &UnsafeValue{
		t:     t,
		model: model,
		addr:  addr,
	}
}

func (u *UnsafeValue) SetColumns(rows *sql.Rows) error {
	if !rows.Next() {
		return err2.ErrNoRows
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(columns) > len(u.model.FieldMap) {
		return err2.ErrTooManyReturnedColumns
	}

	vals := make([]any, 0, len(columns))
	for _, column := range columns {
		f := u.model.ColumnMap[column]

		fdVal := reflect.NewAt(f.Typ, unsafe.Pointer(uintptr(u.addr)+f.Offset))
		vals = append(vals, fdVal.Interface())
	}

	return rows.Scan(vals...)
}
