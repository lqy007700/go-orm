package go_orm

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-orm/internal/valuer"
	"log"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	db := memoryDB()
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
	db, mock, err := sqlmock.New()
	if err != nil {
		return
	}
	defer db.Close()

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
				res.AddRow([]byte("1"), []byte("Liu"), []byte("18"), []byte("Quan"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Liu",
				Age:       18,
				LastName:  &sql.NullString{String: "Quan", Valid: true},
			},
		},
	}

	openDB, err := OpenDB(db)
	if err != nil {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := mock.ExpectQuery(tt.query)
			if tt.mockErr != nil {
				exp.WillReturnError(tt.mockErr)
			} else {
				exp.WillReturnRows(tt.mockRows)
			}

			get, err := NewSelector[TestModel](openDB).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, get)
		})
	}
}

func memoryDB() *DB {
	orm, err := Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	return orm
}

func BenchmarkQuerier_Get(b *testing.B) {
	db := memoryDB()
	res, err := db.db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)"+
		"VALUES (?,?,?,?)", 13, "Deng", 18, "Ming")

	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	log.Println(affected)
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}
	b.Run("unsafe", func(b *testing.B) {
		db.valCreator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.valCreator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
