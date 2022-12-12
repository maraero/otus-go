package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upInit, downInit)
}

func upInit(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE eents (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			date_start TIMESTAMP NOT NULL,
			date_end TIMESTAMP NOT NULL, 
			description TEXT NULL,
			user_id TEXT NOT NULL,
			date_notification TIMESTAMP NULL
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func downInit(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events")
	if err != nil {
		return err
	}
	return nil
}
