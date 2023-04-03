package go_orm

import (
	"context"
	"errors"
	err2 "go-orm/internal/err"
	model2 "go-orm/internal/model"
	"strings"
)

type Selector[T any] struct {
	builder
	sb      strings.Builder
	Args    []any
	table   string
	where   []Predicate
	having  []Predicate
	columns []Selectable
	groupBy []Column
	orderBy []OrderBy
	limit   int32
	offset  int32

	model *model2.Model
	//db    *DB
	sess Session
	core
}

type Selectable interface {
	selectable()
}

func NewSelector[T any](sess Session) *Selector[T] {
	return &Selector[T]{
		sess: sess,
		core: sess.getCore(),
	}
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}

func (s *Selector[T]) OrderBy(ps ...OrderBy) *Selector[T] {
	s.orderBy = ps
	return s
}

func (s *Selector[T]) GroupBy(ps ...Column) *Selector[T] {
	s.groupBy = ps
	return s
}

func (s *Selector[T]) Limit(l int32) *Selector[T] {
	s.limit = l
	return s
}

func (s *Selector[T]) Offset(o int32) *Selector[T] {
	s.offset = o
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   = new(T)
		err error
	)
	s.model, err = s.r.Get(t)
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")

	s.buildTableName()

	if err = s.buildWhere(); err != nil {
		return nil, err
	}

	if err = s.buildGroupBy(); err != nil {
		return nil, err
	}

	if err = s.buildHaving(); err != nil {
		return nil, err
	}

	if err = s.buildOrderBy(); err != nil {
		return nil, err
	}

	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
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
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		build, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}

		rows, err := s.sess.queryContext(ctx, build.SQL, s.Args...)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		t := new(T)

		val := s.valCreator(t, s.model)
		return &QueryResult{
			Result: t,
			Err:    val.SetColumns(rows),
		}
	}

	for i := len(s.ms) - 1; i >= 0; i-- {
		root = s.ms[i](root)
	}

	res := root(ctx, &QueryContext{
		Type:    "SELECT",
		Builder: s,
	})

	if res.Err != nil {
		return nil, res.Err
	}

	t, ok := res.Result.(*T)
	if !ok {
		return nil, errors.New("类型错误")
	}

	return t, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	return nil, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
	} else {
		for i, column := range s.columns {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			switch c := column.(type) {
			case Column:
				if err := s.buildColumn(c.name, c.alias); err != nil {
					return err
				}
			case Aggregate:
				if err := s.buildAggregate(c); err != nil {
					return err
				}
			case RawExpr:
				s.sb.WriteString(c.raw)
				if len(c.args) > 0 {
					s.Args = append(s.Args, c.args)
				}
			}
		}
	}
	return nil
}

func (s *Selector[T]) buildColumn(name, alias string) error {
	s.sb.WriteByte('`')
	fd, ok := s.model.FieldMap[name]
	if !ok {
		return err2.NewErrUnknownColumn(name)
	}
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')

	s.buildAs(alias)
	return nil
}

func (s *Selector[T]) buildAs(alias string) {
	if alias != "" {
		s.sb.WriteString(" as ")
		s.sb.WriteByte('`')
		s.sb.WriteString(alias)
		s.sb.WriteByte('`')
	}
}

func (s *Selector[T]) buildAggregate(c Aggregate) error {
	s.sb.WriteString(c.fn)
	s.sb.WriteByte('(')

	fd, ok := s.model.FieldMap[c.arg]
	if !ok {
		return err2.NewErrUnknownColumn(c.arg)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')
	s.sb.WriteByte(')')

	s.buildAs(c.alias)

	return nil
}

func (s *Selector[T]) buildTableName() {
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
}

func (s *Selector[T]) buildWhere() error {
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		return s.buildPredicates(s.where)
	}
	return nil
}

func (s *Selector[T]) buildGroupBy() error {
	if len(s.groupBy) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for i, c := range s.groupBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}

			if err := s.buildColumn(c.name, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Selector[T]) buildHaving() error {
	if len(s.having) > 0 {
		s.sb.WriteString(" HAVING ")
		return s.buildPredicates(s.having)
	}
	return nil
}

func (s *Selector[T]) addArgs(args ...any) {
	if s.Args == nil {
		s.Args = make([]any, 0, 8)
	}
	s.Args = append(s.Args, args...)
}

func (s *Selector[T]) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	return s.buildExpression(p)
}

func (s *Selector[T]) buildOrderBy() error {
	if len(s.orderBy) > 0 {
		s.sb.WriteString(" ORDER BY ")

		for i, by := range s.orderBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}

			s.sb.WriteByte('`')
			s.sb.WriteString(by.col)
			s.sb.WriteByte('`')
			s.sb.WriteString(" " + by.order)
		}
	}
	return nil
}

type OrderBy struct {
	col   string
	order string
}

func Asc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "ASC",
	}
}

func Desc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "DESC",
	}
}
