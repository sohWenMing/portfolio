package dbinfra

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConfig interface {
	DBString() string
}

func InitDB(dbc DBConfig) (dbconn *sql.DB, err error) {
	db, err := sql.Open("postgres", dbc.DBString())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
