package go_orm

import (
	"errors"
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	sb   strings.Builder
	db   *DB
	vals []*T
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

	for i2, column := range m.Columns {
		if i2 > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(column.ColName)
		i.sb.WriteByte('`')
	}
	i.sb.WriteByte(')')

	i.sb.WriteString(" VALUES(")

	args := make([]any, 0, len(i.vals)*len(m.Columns))
	//val := i.vals[0]

	for i3, val := range i.vals {
		of := reflect.ValueOf(val).Elem()
		if i3 > 0 {
			i.sb.WriteString("),(")
		}

		for i2, c := range m.Columns {
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
