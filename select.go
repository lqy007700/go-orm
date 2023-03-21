package go_orm

import (
	"context"
	"errors"
	model2 "go-orm/internal/model"
	"strings"
)

type Selector[T any] struct {
	sb    strings.Builder
	Args  []any
	table string
	where []Predicate
	model *model2.Model

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   = new(T)
		err error
	)
	s.model, err = s.db.r.Get(t)
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT * FROM ")
	s.sb.WriteByte('`')

	if s.table == "" {
		s.sb.WriteString(s.model.TableName)
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
		Args: s.Args,
	}, nil
}

func (s *Selector[T]) buildExpression(expression Expression) error {
	switch expr := expression.(type) {
	case nil:
		return nil
	case Column:
		s.sb.WriteByte('`')
		f, ok := s.model.FieldMap[expr.name]
		if !ok {
			return errors.New("字段不存在")
		}
		s.sb.WriteString(f.ColName)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteByte('?')
		if s.Args == nil {
			s.Args = make([]any, 0, 8)
		}
		s.Args = append(s.Args, expr.val)
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
		if expr.op.String() != "NOT" {
			s.sb.WriteByte(' ')
		}

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

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	build, err := s.Build()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.db.QueryContext(ctx, build.SQL, s.Args...)
	if err != nil {
		return nil, err
	}
	t := new(T)

	val := s.db.valCreator(t, s.model)
	return t, val.SetColumns(rows)
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//var db *sql.DB
	//build, err := s.Build()
	//if err != nil {
	//	return nil, err
	//}
	//
	//queryContext, err := db.QueryContext(ctx, build.SQL, build.args)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}
