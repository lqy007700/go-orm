package go_orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
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
	db := memoryDB(t)
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
				Args: []any{12},
			},
			wantErr: nil,
		},
		{
			name: "where all",
			s:    NewSelector[TestModel](db).Where(C("Age").GT(18), C("Age").LT(35)),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
			wantErr: nil,
		},
		{
			name: "not",
			s:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE NOT (`age` > ?);",
				Args: []any{18},
			},
			wantErr: nil,
		},
		{
			name: "and",
			s:    NewSelector[TestModel](db).Where(C("Age").EQ(12), C("FirstName").EQ("liu")),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`first_name` = ?);",
				Args: []any{12, "liu"},
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

func TestSelector_Get(t *testing.T) {
	mockDb, _, err := sqlmock.New()
	if err != nil {
		return
	}

	tests := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
		{
			// 查询返回错误
			name:    "query error",
			mockErr: errors.New("invalid query"),
			wantErr: errors.New("invalid query"),
			query:   "SELECT .*",
		},
	}

	db, err := OpenDB(mockDb)
	if err != nil {
		return
	}

	//for _, tc := range tests {
	//	exp := mock.ExpectQuery(tc.query)
	//	if tc.mockErr != nil {
	//		exp.WillReturnError(tc.mockErr)
	//	} else {
	//		exp.WillReturnRows(tc.mockRows)
	//	}
	//}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, get)
		})
	}
}

func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}
