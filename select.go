package go_orm

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	sb    strings.Builder
	table string
	where []Predicate

	args []any
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

//func (s *Selector[T]) Where(where string, args ...any) *Selector[T] {
//	return s
//}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb.WriteString("SELECT * FROM ")
	s.sb.WriteByte('`')

	if s.table == "" {
		var t T
		// 使用结构体名做为表名
		of := reflect.TypeOf(t)
		name := of.Name()
		s.sb.WriteString(name)
	} else {
		// 处理 db.table_name 的情况
		segs := strings.SplitN(s.table, ".", 2)
		if len(segs) == 2 {
			s.sb.WriteString(segs[0])
			s.sb.WriteByte('`')
			s.sb.WriteByte('.')
			s.sb.WriteByte('`')
			s.sb.WriteString(segs[1])
		} else {
			s.sb.WriteString(s.table)
		}
	}
	s.sb.WriteByte('`')

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		pred := s.where[0]
		for i := 1; i < len(s.where); i++ {
			pred = pred.And(s.where[i])
		}
		err := s.buildExpression(pred)
		if err != nil {
			return nil, err
		}
	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		args: s.args,
	}, nil
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) buildExpression(expression Expression) error {
	switch expr := expression.(type) {
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(expr.name)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteByte('?')
		if s.args == nil {
			s.args = make([]any, 0, 8)
		}
		s.args = append(s.args, expr.val)
	case Predicate:
		_, ok := expr.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		err := s.buildExpression(expr.left)
		if err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(expr.op.String())
		s.sb.WriteByte(' ')

		_, ok = expr.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}

		err = s.buildExpression(expr.right)
		if ok {
			s.sb.WriteByte(')')
		}
		return err
	default:
		return errors.New("不支持该表达式")
	}
	return nil
}
