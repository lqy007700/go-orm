package go_orm

import (
	"errors"
	err2 "go-orm/internal/err"
	"go-orm/internal/model"
	"reflect"
)

type Inserter[T any] struct {
	builder
	db   *DB
	vals []*T
	cols []string

	onDuplicate *OnDuplicateKey
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: builder{
			dialect: db.dialect,
		},
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
	i.m = m
	i.builder.quote(m.TableName)
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

		i.builder.quote(field.ColName)
	}

	i.sb.WriteString(")VALUES(")
	i.args = make([]any, 0, len(i.vals)*len(m.Columns))
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
			i.args = append(i.args, of.FieldByIndex(c.Index).Interface())
		}
	}
	i.sb.WriteString(")")

	// 构造 onDuplicate
	if i.onDuplicate != nil {
		err = i.dialect.buildDuplicateKey(&i.builder, i.onDuplicate)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteString(";")

	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}

func (o *OnDuplicateKeyBuilder[T]) update(assign ...Assignable) *Inserter[T] {
	o.i.onDuplicate = &OnDuplicateKey{
		assigns: assign,
	}
	return o.i
}

type OnDuplicateKey struct {
	assigns []Assignable
}
