package postgres

import (
	"context"
	"database/sql"

	"ops-storage/internal/server/logger"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) Database {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.MainLog.Panicf("Postgres open error, err: ", err.Error())
	}
	res := Database{
		db: db,
	}

	if err = res.pingBackoff(); err != nil {
		logger.MainLog.Panicf("Postgres connection error, err: ", err.Error())
	}

	res.prep()

	return res
}

// TODO parse a taxonomy
func (d Database) prep() {
	query := `
		CREATE TABLE IF NOT EXISTS gauge (
			id SERIAL PRIMARY KEY, 
			name varchar(32) UNIQUE, 
			value double precision NOT NULL,
			updated_at timestamp DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS counter (
			id SERIAL PRIMARY KEY, 
			name varchar(32) UNIQUE, 
			value integer NOT NULL,
			updated_at timestamp DEFAULT NOW()
		);
	`
	_, err := d.db.ExecContext(context.Background(), query)
	if err != nil {
		logger.MainLog.Panicf("Postgres tables creation error, err: ", err.Error())
	}
}
