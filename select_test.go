package go_orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "from",
			s:    NewSelector[TestModel]().From("test_model_tab"),
			want: &Query{
				SQL: "SELECT * FROM `test_model_tab`;",
			},
			wantErr: nil,
		},
		{
			name: "not from",
			s:    NewSelector[TestModel](),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
			wantErr: nil,
		},
		{
			name: "with db",
			s:    NewSelector[TestModel]().From("test_db.test_model"),
			want: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
			wantErr: nil,
		},
		{
			name: "where",
			s:    NewSelector[TestModel]().Where(C("id").EQ(12)),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE `id`=?;",
				args: []any{12},
			},
			wantErr: nil,
		},
		{
			name: "where all",
			s:    NewSelector[TestModel]().Where(C("age").GT(18), C("age").LT(35)),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`age` > ?) AND (`age` < ?);",
				args: []any{18, 35},
			},
			wantErr: nil,
		},
		{
			name: "not",
			s:    NewSelector[TestModel]().Where(Not(C("age").GT(18))),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE NOT (`age` > ?);",
				args: []any{18},
			},
			wantErr: nil,
		},
		{
			name: "and",
			s:    NewSelector[TestModel]().Where(C("age").EQ(12), C("name").EQ("liu")),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`age` = ?) AND (`name` = ?);",
				args: []any{12, "liu"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
