package model

import (
	"database/sql"
	"echo-web/module/log"
)

var db *sql.DB

// DB n
func DB() *sql.DB {
	if db == nil {
		log.DebugPrint("Model NewDB")
		newDb, err := NewDB()
		if err != nil {
			panic(err)
		}
		db = newDb
	}

	return db
}
