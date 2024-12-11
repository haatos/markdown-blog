package data

import (
	"database/sql"
	"log"

	assets "github.com/haatos/markdown-blog"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) {
	/*
		run goose migrations to the newest version
	*/
	goose.SetBaseFS(assets.MigrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}
}
