package db

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DBConnection(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("Error opening DB connection => %s", err.Error())

	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging DB => %s", err.Error())
	}

	return db, nil
}
