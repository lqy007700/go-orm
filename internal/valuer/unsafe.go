package valuer

import (
	"database/sql"
)

type unsafeValue struct {
}

func (u unsafeValue) SetColumns(row *sql.Rows) error {
	//TODO implement me
	panic("implement me")
}
