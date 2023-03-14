package go_orm

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

func Test_registry_get(t *testing.T) {
	type fields struct {
		models map[reflect.Type]*model
		lock   sync.RWMutex
	}

	tests := []struct {
		name    string
		fields  fields
		val     any
		want    *model
		wantErr assert.ErrorAssertionFunc
	}{
		// 标签相关测试用例
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			want: &model{
				tableName: "column_tag",
				fieldMap: map[string]*field{
					"ID": {
						colName: "id",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//r := &registry{
			//	models: tt.fields.models,
			//	lock:   sync.RWMutex{},
			//}
		})
	}
}

func Test_underscoreName(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "123",
			s:    "IdA",
			want: "id_a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, underscoreName(tt.s), "underscoreName(%v)", tt.s)
		})
	}
}
