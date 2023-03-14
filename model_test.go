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
	}

	r := &registry{
		models: map[reflect.Type]*model{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := r.parseModel(tt.input)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})
	}
}
