package model

import (
	"github.com/stretchr/testify/assert"
	"go-orm"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		opts    []Opt
		want    *model
		wantErr error
	}{
		{
			name:  "ptr",
			input: &go_orm.TestModel{},
			want: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id":        {ColName: "id"},
					"FirstName": {ColName: "first_name"},
					"Age":       {ColName: "age"},
					"LastName":  {ColName: "last_name"},
				},
			},
			wantErr: nil,
		},
		{
			name:  "struct",
			input: go_orm.TestModel{},
			want: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id":        {ColName: "id"},
					"FirstName": {ColName: "first_name"},
					"Age":       {ColName: "age"},
					"LastName":  {ColName: "last_name"},
				},
			},
			wantErr: nil,
		},
		{
			name:    "nil",
			input:   nil,
			want:    nil,
			wantErr: nil,
		},
		{
			name: "column tag",
			input: func() any {
				type ColumnTag struct {
					ID uint64 `orm:"column=ids"`
				}
				return &ColumnTag{}
			}(),
			want: &model{
				tableName: "column_tag",
				fieldMap: map[string]*field{
					"ID": {
						ColName: "ids",
					},
				},
			},
		},
		{
			name:  "with table name ",
			input: go_orm.TestModel{},
			opts:  []Opt{WithTableName("a"), WithColumnName("Id", "uid")},
			want: &model{
				tableName: "a",
				fieldMap: map[string]*field{
					"Id":        {ColName: "uid"},
					"FirstName": {ColName: "first_name"},
					"Age":       {ColName: "age"},
					"LastName":  {ColName: "last_name"},
				},
			},
		},
	}

	r := &go_orm.registry{
		models: map[reflect.Type]*model{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.Register(tt.input, tt.opts...)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})
	}
}
