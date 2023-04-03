package querylog

import (
	"context"
	"fmt"
	go_orm "go-orm"
	"time"
)

type MiddlewareBuilder struct {
	threshold int64
	logFunc   func(sql string, args ...any)
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args ...any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m *MiddlewareBuilder) Build() go_orm.MiddleWare {
	return func(next go_orm.Handler) go_orm.Handler {
		return func(ctx context.Context, qc *go_orm.QueryContext) *go_orm.QueryResult {
			start := time.Now()
			build, err := qc.Builder.Build()
			if err != nil {
				return &go_orm.QueryResult{
					Result: nil,
					Err:    err,
				}
			}
			m.logFunc(build.SQL, build.Args...)

			defer func() {
				// 慢查询监控
				duration := time.Now().Sub(start)
				if m.threshold > 0 && duration.Microseconds() > m.threshold {
					fmt.Println(123)
				}
			}()

			return next(ctx, qc)
		}
	}
}
