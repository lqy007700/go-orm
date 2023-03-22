package go_orm

// Expression 标记接口
// expr 并不具有实际意义
type Expression interface {
	expr()
}

type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectable() {}

func Raw(raw string, args ...any) RawExpr {
	return RawExpr{
		raw:  raw,
		args: args,
	}
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{}
}
