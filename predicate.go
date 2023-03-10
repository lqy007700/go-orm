package go_orm

type op string

const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opNOT = "NOT"
	opAND = "AND"
	opOR  = "OR"
)

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (o op) String() string {
	return string(o)
}

func (p Predicate) expr() {}

type Column struct {
	name string
}

func (c Column) expr() {}

type Value struct {
	val any
}

func (v Value) expr() {}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) EQ(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: Value{val: val},
	}
}

func (c Column) GT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: Value{val: val},
	}
}

func (c Column) LT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: Value{val: val},
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (p Predicate) And(p2 Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: p2,
	}
}

func (p Predicate) Or(p2 Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: p2,
	}
}
