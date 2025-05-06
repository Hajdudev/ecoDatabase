package store

import (
	"database/sql"
	"fmt"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

type DatabaseStore interface {
	getRoutes(id string)
}

func (pg *PostgresStore) getRoutes(id string) {
	fmt.Println("Found something ")
}
