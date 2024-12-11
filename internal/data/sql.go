package data

import (
	"database/sql"
	"log"
	"runtime"

	"github.com/haatos/markdown-blog/internal"
)

func InitDatabase(readonly bool) *sql.DB {
	db, err := sql.Open("sqlite3", internal.Settings.SQLiteDbString(readonly))
	if err != nil {
		log.Fatal(err)
	}

	if readonly {
		db.SetMaxOpenConns(max(4, runtime.NumCPU()))
	} else {
		if _, err := db.Exec("PRAGMA temp_store=memory;"); err != nil {
			log.Fatal(err)
		}
		db.SetMaxOpenConns(1)
	}

	return db
}
