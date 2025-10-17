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
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
  comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(256) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	if install {
		if _, err = conn.Exec(schema); err != nil {
			return fmt.Errorf("schema init failed: %w", err)
		}
	}
	db = conn

	return nil
}

func Close() error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db connection")
	}
	return nil
}
