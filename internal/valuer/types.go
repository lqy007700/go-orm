package valuer

import "database/sql"

// Value reflect 和 unsafe 的抽象
type Value interface {
	SetColumns(row *sql.Rows) error
}
