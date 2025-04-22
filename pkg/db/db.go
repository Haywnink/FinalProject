package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_date ON scheduler(date);
`

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	_, err := os.Stat(path)
	isNew := os.IsNotExist(err)
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if isNew {
		if _, err := conn.Exec(schema); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to create schema: %v", err)
		}
	}
	return &DB{conn: conn}, nil
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) Close() error {
	return db.conn.Close()
}
