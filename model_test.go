package go_orm

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		opts    []ModelOpt
		want    *model
		wantErr error
	}{
		{
			name:  "ptr",
			input: &TestModel{},
			want: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id":        {colName: "id"},
					"FirstName": {colName: "first_name"},
					"Age":       {colName: "age"},
					"LastName":  {colName: "last_name"},
				},
			},
			wantErr: nil,
		},
		{
			name:  "struct",
			input: TestModel{},
			want: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id":        {colName: "id"},
					"FirstName": {colName: "first_name"},
					"Age":       {colName: "age"},
					"LastName":  {colName: "last_name"},
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
						colName: "ids",
					},
				},
			},
		},
		{
			name:  "with table name ",
			input: TestModel{},
			opts:  []ModelOpt{ModelWithTableName("a"), ModelWithColumnName("Id", "uid")},
			want: &model{
				tableName: "a",
				fieldMap: map[string]*field{
					"Id":        {colName: "uid"},
					"FirstName": {colName: "first_name"},
					"Age":       {colName: "age"},
					"LastName":  {colName: "last_name"},
				},
			},
		},
	}

	r := &registry{
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
