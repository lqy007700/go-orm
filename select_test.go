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

	db, _ := NewDB()

	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "from",
			s:    NewSelector[TestModel](db).From("test_model_tab"),
			want: &Query{
				SQL: "SELECT * FROM `test_model_tab`;",
			},
			wantErr: nil,
		},
		{
			name: "not from",
			s:    NewSelector[TestModel](db),
			want: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
			wantErr: nil,
		},
		{
			name: "with db",
			s:    NewSelector[TestModel](db).From("test_db.test_model"),
			want: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
			wantErr: nil,
		},
		{
			name: "where",
			s:    NewSelector[TestModel](db).Where(C("Id").EQ(12)),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				args: []any{12},
			},
			wantErr: nil,
		},
		{
			name: "where all",
			s:    NewSelector[TestModel](db).Where(C("Age").GT(18), C("Age").LT(35)),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				args: []any{18, 35},
			},
			wantErr: nil,
		},
		{
			name: "not",
			s:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE NOT (`age` > ?);",
				args: []any{18},
			},
			wantErr: nil,
		},
		{
			name: "and",
			s:    NewSelector[TestModel](db).Where(C("Age").EQ(12), C("FirstName").EQ("liu")),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
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
