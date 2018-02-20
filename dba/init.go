// Package dba  provides database accessor
package dba

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var ErrAlreadyExists = errors.New("already exists")

func init() {
	dbpath := "./db/sqlite3.db"
	if _, err := os.Stat(dbpath); err != nil {
		if err == os.ErrNotExist {
			os.Create(dbpath)
		}
	}

	var err error
	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
}
