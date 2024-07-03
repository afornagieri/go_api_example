package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqlCli struct {
	Conn *sql.DB
}

func NewSqlCli() (*SqlCli, error) {
	conn, err := sql.Open("sqlite3", "./internal/infra/database/items.db")
	ensureTableExists(conn)
	if err != nil {
		return nil, err
	}
	return &SqlCli{Conn: conn}, nil
}

func ensureTableExists(db *sql.DB) error {
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS items (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL,
					price REAL NOT NULL,
					description TEXT NOT NULL
			)
	`)
	if err != nil {
		return err
	}
	return nil
}
