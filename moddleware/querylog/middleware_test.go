package querylog

import (
	"context"
	"fmt"
	go_orm "go-orm"
	"testing"
)

type TestModel struct {
}

func TestMiddlewareBuilder_Build(t *testing.T) {
	build := &MiddlewareBuilder{
		threshold: 9000,
	}
	orm, err := go_orm.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8mb4", go_orm.DBWithMiddleware(
		build.LogFunc(func(sql string, args ...any) {
			fmt.Println(sql)
		}).Build()))
	if err != nil {
		panic(err)
	}

	get, err := go_orm.NewSelector[TestModel](orm).Get(context.Background())

	if err != nil {
		return
	}

	fmt.Println(get)
}
