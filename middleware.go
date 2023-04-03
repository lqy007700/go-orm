package go_orm

import "context"

type QueryContext struct {
	Type    string
	Builder QueryBuilder
}

type QueryResult struct {
	Result any

	Err error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type MiddleWare func(next Handler) Handler
