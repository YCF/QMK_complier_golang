package model

import (
	"Goose/conf"
	"database/sql"
)

// NewDB n
func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", conf.GetOption("system", "db"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
