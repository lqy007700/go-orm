package model

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opt ...Opt) (*Model, error)
}

type Registrys struct {
	Models map[reflect.Type]*Model
	lock   sync.Map
}

func (r *Registrys) Get(val any) (*Model, error) {
	of := reflect.TypeOf(val)

	value, ok := r.lock.Load(of)
	if ok {
		return value.(*Model), nil
	}
	return r.Register(val)
}

func (r *Registrys) Register(val any, opts ...Opt) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(m)
	}

	typ := reflect.TypeOf(val)
	r.lock.Store(typ, m)
	return m, nil
}

func (r *Registrys) parseModel(val any) (*Model, error) {
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
	columnMap := make(map[string]*field)
	for i := 0; i < fieldCnt; i++ {
		fd := of.Field(i)

		var colName string
		var ok bool
		if fd.Tag == "" {
			colName = underscoreName(fd.Name)
		} else {
			tags, err := parseTag(fd.Tag)
			if err != nil {
				return nil, err
			}
			colName, ok = tags["column"]
			if !ok || colName == "" {
				colName = underscoreName(fd.Name)
			}
		}

		fieldV := &field{
			GoName:  fd.Name,
			ColName: colName,
			Typ:     fd.Type,
			Offset:  fd.Offset,
		}

		fieldMap[fd.Name] = fieldV
		columnMap[colName] = fieldV
	}

	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(of.Name())
	}

	return &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
	}, nil
}

func parseTag(tag reflect.StructTag) (map[string]string, error) {
	tagv := tag.Get("orm")

	kvs := strings.Split(tagv, ",")

	res := make(map[string]string, 1)
	for _, kv := range kvs {
		n := strings.Split(kv, "=")

		k := n[0]
		var v string
		if len(n) > 1 {
			v = n[1]
		}
		res[k] = v
	}
	return res, nil
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
