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

var Conn *sql.DB

func Init(defaultPath string) error {
	path := os.Getenv("TODO_DBFILE")
	if path == "" {
		path = defaultPath
	}
	_, err := os.Stat(path)
	first := os.IsNotExist(err)
	Conn, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %v", err)
	}
	if first {
		if _, err = Conn.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы БД: %v", err)
		}
	}
	return nil
}
