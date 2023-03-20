package go_orm

import "reflect"

type ModelOpt func(m *model)

type model struct {
	tableName string
	fieldMap  map[string]*field

	// 列名-字段名
	columnMap map[string]*field
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
	// 字段名
	goName string

	// 列名
	colName string

	typ reflect.Type
}

type TableName interface {
	TableName() string
}
