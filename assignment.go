package go_orm

type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    any
}

func (a Assignment) assign() {}

func Assign(col string, val any) Assignment {
	return Assignment{
		column: col,
		val:    val,
	}
}

type CC struct {
	column string
	val    any
}
