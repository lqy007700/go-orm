package go_orm

type Aggregate struct {
	arg   string
	fn    string
	alias string
}

func (a Aggregate) selectable() {}

func Avg(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "AVG",
	}
}

func Min(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MIN",
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MAX",
	}
}

func (a Aggregate) As(as string) Aggregate {
	a.alias = as
	return a
}
