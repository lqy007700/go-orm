package go_orm

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type registry struct {
	models map[reflect.Type]*model
	lock   sync.RWMutex
}

func (r *registry) get(val any) (*model, error) {
	r.lock.RLock()
	of := reflect.TypeOf(val)
	m, ok := r.models[of]
	r.lock.RUnlock()
	if ok {
		return m, nil
	}

	r.lock.Lock()
	defer r.lock.Lock()
	m, ok = r.models[of]
	r.lock.RUnlock()
	if ok {
		return m, nil
	}

	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models[of] = m
	return m, nil
}

func (r *registry) parseModel(val any) (*model, error) {
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

		fieldMap[fd.Name] = &field{
			colName: colName,
		}
	}

	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(of.Name())
	}

	return &model{
		tableName: tableName,
		fieldMap:  fieldMap,
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
