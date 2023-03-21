package model

import "reflect"

type ModelOpt func(m *Model)

type Model struct {
	TableName string
	FieldMap  map[string]*field

	// 列名-字段名
	ColumnMap map[string]*field
}

func ModelWithTableName(name string) ModelOpt {
	return func(m *Model) {
		m.TableName = name
	}
}

func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) {
		m.FieldMap[field].ColName = colName
	}
}

func ModelWithColumn(field string, col *field) ModelOpt {
	return func(m *Model) {
		m.FieldMap[field] = col
	}
}

type field struct {
	// 字段名
	GoName string

	// 列名
	ColName string

	Typ reflect.Type

	// 表达相对量的概念
	// 偏移量
	offset uintptr
}

type TableName interface {
	TableName() string
}
