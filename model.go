package go_orm

type ModelOpt func(m *model)

type model struct {
	tableName string
	fieldMap  map[string]*field
}

func ModelWithTableName(name string) ModelOpt {
	return func(m *model) {
		m.tableName = name
	}
}

func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *model) {
		m.fieldMap[field].colName = colName
	}
}

func ModelWithColumn(field string, col *field) ModelOpt {
	return func(m *model) {
		m.fieldMap[field] = col
	}
}

type field struct {
	colName string
}

type TableName interface {
	TableName() string
}
