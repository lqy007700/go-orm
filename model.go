package go_orm

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}
