package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to Database...")
	return db, nil
}
