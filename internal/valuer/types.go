package valuer

import (
	"database/sql"
	"go-orm/internal/model"
)

// Value reflect 和 unsafe 的抽象
type Value interface {
	SetColumns(row *sql.Rows) error
}

type Creator func(t any, model *model.Model) Value
