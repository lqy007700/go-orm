package model

import "reflect"

type Opt func(m *Model)

type Model struct {
	TableName string
	FieldMap  map[string]*field
	Columns   []*field

	// 列名-字段名
	ColumnMap map[string]*field
}

func WithTableName(name string) Opt {
	return func(m *Model) {
		m.TableName = name
	}
}

func WithColumnName(field string, colName string) Opt {
	return func(m *Model) {
		m.FieldMap[field].ColName = colName
	}
}

func WithColumn(field string, col *field) Opt {
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
	Offset uintptr

	Index []int
}

type TableName interface {
	TableName() string
}
