package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetDB() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("database url not specified")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic("cant connect to database" + err.Error())
	}

	return db
}

func CreateTables(db *sql.DB) error {
	data, err := os.ReadFile("./sql/InitDB.sql")
	if err != nil {
		panic("error reading sql file" + err.Error())
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("error creating tables: %s", err.Error())
	}

	return nil
}
