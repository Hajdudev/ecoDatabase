package store

import (
	"context"
	"fmt"
	"os"

	"github.com/Hajdudev/ecoDatabase/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}

type DatabaseStore interface {
	GetUserByID(id string) (*models.User, error)
}

func (pg *PostgresStore) GetUserByID(id string) (*models.User, error) {
	query := "SELECT id, created_at, email, name, image, recent_rides FROM users WHERE id = $1"

	var user models.User
	var recentRidesBytes []string
	err := pg.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Email,
		&user.Name,
		&user.Image,
		&recentRidesBytes,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user by ID: %v\n", err)
		return nil, err
	}

	user.RecentRides = recentRidesBytes

	return &user, nil
}
