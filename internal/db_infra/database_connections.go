package dbinfra

import (
	"database/sql"
	"embed"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type AppDB struct {
	DB *sql.DB
}

func (a *AppDB) Migrate() error {

	goose.SetBaseFS((embedMigrations))
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(a.DB, "migrations"); err != nil {
		return err
	}
	return nil
}

type DBConfig interface {
	DBString() string
}

func InitDB(dbc DBConfig) (appDB *AppDB, err error) {
	db, err := sql.Open("postgres", dbc.DBString())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &AppDB{DB: db}, nil
}
