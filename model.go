package go_orm

import (
	"reflect"
	"unicode"
)

type model struct {
	tableName string
	fieldMap  map[string]*field
}

type field struct {
	colName string
}

func parseModel(val any) (*model, error) {
	if val == nil {
		return nil, nil
	}
	of := reflect.TypeOf(val)
	// 反射解析指针
	for of.Kind() == reflect.Ptr {
		of = of.Elem()
	}

	fieldCnt := of.NumField()

	fieldMap := make(map[string]*field)
	for i := 0; i < fieldCnt; i++ {
		fd := of.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}

	return &model{
		tableName: underscoreName(of.Name()),
		fieldMap:  fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}
