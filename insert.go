package go_orm

import (
	"errors"
	err2 "go-orm/internal/err"
	"go-orm/internal/model"
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	sb   strings.Builder
	db   *DB
	vals []*T
	cols []string
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.vals = vals
	return i
}

func (i *Inserter[T]) Columns(vals ...string) *Inserter[T] {
	i.cols = vals
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.vals) <= 0 {
		return nil, errors.New("err")
	}

	i.sb.WriteString("INSERT INTO ")

	m, err := i.db.r.Get(i.vals[0])
	if err != nil {
		return nil, err
	}
	i.sb.WriteByte('`')
	i.sb.WriteString(m.TableName)
	i.sb.WriteByte('`')
	i.sb.WriteByte('(')

	fields := m.Columns
	if len(i.cols) != 0 {
		fields = make([]*model.Field, 0, len(i.cols))

		for _, col := range i.cols {
			fd, ok := m.FieldMap[col]
			if !ok {
				return nil, err2.NewErrUnknownColumn(col)
			}
			fields = append(fields, fd)
		}
	}

	for i2, field := range fields {
		if i2 > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('`')
		i.sb.WriteString(field.ColName)
		i.sb.WriteByte('`')
	}

	i.sb.WriteString(")VALUES(")
	args := make([]any, 0, len(i.vals)*len(m.Columns))
	for i3, val := range i.vals {
		of := reflect.ValueOf(val).Elem()
		if i3 > 0 {
			i.sb.WriteString("),(")
		}

		for i2, c := range fields {
			if i2 > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			args = append(args, of.FieldByIndex(c.Index).Interface())
		}
	}

	i.sb.WriteString(");")
	return &Query{
		SQL:  i.sb.String(),
		Args: args,
	}, nil
}
