package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func EnsureDB() *sql.DB {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "final_project", "scheduler.db")
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	_, err = os.Stat(filepath.Join(dbFile))

	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			log.Fatal(err)
		}

		query := `CREATE TABLE IF NOT EXISTS "scheduler" (
				"id" INTEGER NOT NULL UNIQUE,
				"date" DATE NOT NULL DEFAULT "",
				"title" VARCHAR(128) NOT NULL DEFAULT "",
				"comment" TEXT NOT NULL DEFAULT "",
				"repeat" VARCHAR(128) NOT NULL DEFAULT "",
				PRIMARY KEY("id")
			);
		
			CREATE INDEX IF NOT EXISTS "scheduler_index_date"
			ON "scheduler" ("date");`

		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
		return db
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
