package go_orm

import (
	err2 "go-orm/internal/err"
)

type Dialect interface {
	quoter() byte
	buildDuplicateKey(b *builder, key *OnDuplicateKey) error
}

// 标准sql
type standardSQL struct {
}

type mysqlDialect struct {
	standardSQL
}

func (d *mysqlDialect) quoter() byte {
	return '`'
}

func (d *mysqlDialect) buildDuplicateKey(bu *builder, odk *OnDuplicateKey) error {
	bu.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i2, assign := range odk.assigns {
		if i2 > 0 {
			bu.sb.WriteByte(',')
		}
		switch a := assign.(type) {
		case Assignment:
			fd, ok := bu.m.FieldMap[a.column]
			if !ok {
				return err2.NewErrUnknownColumn(a.column)
			}

			bu.quote(fd.ColName)
			bu.sb.WriteString("=?")
			bu.args = append(bu.args, a.val)
		case Column:
			fd, ok := bu.m.FieldMap[a.name]
			if !ok {
				return err2.NewErrUnknownColumn(a.name)
			}
			bu.quote(fd.ColName)
			bu.sb.WriteString("=VALUES(")
			bu.quote(fd.ColName)
			bu.sb.WriteByte(')')
		}
	}
	return nil
}

type sqliteDialect struct {
	standardSQL
}
