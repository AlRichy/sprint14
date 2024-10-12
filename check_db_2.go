package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

var store Storage

func NewStorage(db *sql.DB) Storage {
	return Storage{db: db}
}

func EnsureDB() {
	// Task 2

	dbFilePath, err := checkDbPath()

	if err != nil {
		initDB(dbFilePath)
	} else {
		log.Fatal(err)
	}
}

func initDB(dbFilePath string) {
	// Первоначальная инициализация БД
	dbFile, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	store = NewStorage(dbFile)
	initCreateTable()
}

func initCreateTable() {
	query := `CREATE TABLE IF NOT EXISTS "scheduler" (
		"id" INTEGER NOT NULL UNIQUE,
		"date" DATE NOT NULL,
		"title" TEXT,
		"comment" TEXT NOT NULL,
		"repeat" VARCHAR,
		PRIMARY KEY("id")
	);

	CREATE INDEX IF NOT EXISTS "scheduler_index_date"
	ON "scheduler" ("date");`

	res, err := store.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.LastInsertId())
}

func checkDbPath() (string, error) {
	appPath, err := os.Executable()
	appPath = filepath.Dir(appPath)
	if err != nil {
		log.Fatal(err)
	}

	// 2*. Определяем путь к файлу БД через переменную среды
	dbFilePath, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		dbFilePath = filepath.Join(appPath, "scheduler.db")
	}
	_, err = os.Stat(dbFilePath)

	return dbFilePath, err
}
