package go_orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {

	type testCase[T any] struct {
		name    string
		i       QueryBuilder
		want    *Query
		wantErr error
	}
	tests := []testCase[TestModel]{
		{
			name: "单行",
			i: NewInserter[TestModel](memoryDB()).Values(&TestModel{
				Id:        12,
				FirstName: "liu",
				Age:       28,
				LastName:  &sql.NullString{Valid: true, String: "quan"},
			}),
			want: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)VALUES(?,?,?,?);",
				Args: []any{int64(12), "liu", int8(28), &sql.NullString{Valid: true, String: "quan"}},
			},
			wantErr: nil,
		},
		{
			name: "多行",
			i: NewInserter[TestModel](memoryDB()).Values(&TestModel{
				Id:        12,
				FirstName: "liu",
				Age:       28,
				LastName:  &sql.NullString{Valid: true, String: "quan"},
			}, &TestModel{
				Id:        1,
				FirstName: "wang",
				Age:       30,
				LastName:  &sql.NullString{Valid: true, String: "liang"},
			}),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{int64(12), "liu", int8(28), &sql.NullString{Valid: true, String: "quan"},
					int64(1), "wang", int8(30), &sql.NullString{Valid: true, String: "liang"}},
			},
			wantErr: nil,
		},
		{
			name: "指定列",
			i: NewInserter[TestModel](memoryDB()).Columns("FirstName", "Age").Values(&TestModel{
				Id:        12,
				FirstName: "liu",
				Age:       28,
				LastName:  &sql.NullString{Valid: true, String: "quan"},
			}, &TestModel{
				Id:        1,
				FirstName: "wang",
				Age:       30,
				LastName:  &sql.NullString{Valid: true, String: "liang"},
			}),
			want: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`age`)VALUES(?,?),(?,?);",
				Args: []any{"liu", int8(28), "wang", int8(30)},
			},
			wantErr: nil,
		},
		{
			name: "upsert",
			i: NewInserter[TestModel](memoryDB()).Values(&TestModel{
				Id:        12,
				FirstName: "liu",
				Age:       28,
				LastName:  &sql.NullString{Valid: true, String: "quan"},
			}).OnDuplicateKey().update(Assign("Age", 19)),
			want: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{int64(12), "liu", int8(28), &sql.NullString{Valid: true, String: "quan"}, 19},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Build()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}

func TestInserter_Exec(t *testing.T) {
	i := Inserter[TestModel]{}

	exec := i.Exec(context.Background())
	id, err := exec.LastInsertId()
	if err != nil {
		return
	}
	fmt.Println(id)
}
