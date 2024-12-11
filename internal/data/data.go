package data

import (
	"database/sql"
	"fmt"
)

func UpdateQuery(table string, attrs ...string) string {
	q := "update " + table + " set "
	for i := range attrs {
		if i == 0 {
			q += attrs[i] + " = " + fmt.Sprintf("$%d", i+1)
		} else {
			q += ", " + attrs[i] + " = " + fmt.Sprintf("$%d", i+1)
		}
	}
	return q
}

func UpdateIDQuery(table string, attrs ...string) string {
	return fmt.Sprintf("%s where id = $%d", UpdateQuery(table, attrs...), len(attrs)+1)
}

func ValuePlaceholders(count int) string {
	q := ""
	for i := range count {
		if i == 0 {
			q += fmt.Sprintf("($%d", i+1)
		} else {
			q += fmt.Sprintf(", $%d", i+1)
		}
	}
	q += ")"
	return q
}

func WithTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	if err := fn(txn); err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}
