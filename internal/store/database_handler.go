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
	GetRoutesById(firstID string, secondID string) ([]string, error)
	GetStopInfo(stopID string) (models.Stop, error)
}

func (pg *PostgresStore) GetStopInfo(stopID string) (models.Stop, error) {
	query := `
		SELECT stop_id, stop_code, stop_name, stop_desc, stop_lat, stop_lon
		FROM stops
		WHERE stop_id = $1
	`

	var stop models.Stop

	err := pg.db.QueryRow(context.Background(), query, stopID).Scan(
		&stop.StopID,
		&stop.StopCode,
		&stop.StopName,
		&stop.StopDesc,
		&stop.StopLat, &stop.StopLon,
	)
	if err != nil {
		return models.Stop{}, fmt.Errorf("error querying stop info: %w", err)
	}

	fmt.Printf("this is the stop %+v", stop)
	return stop, nil
}

func (pg *PostgresStore) GetRoutesById(firstID string, secondID string) ([]string, error) {
	query := `
		SELECT trip_id
		FROM stop_times
		WHERE stop_id IN ($1, $2)
		GROUP BY trip_id
		HAVING COUNT(DISTINCT stop_id) = 2
	`

	rows, err := pg.db.Query(context.Background(), query, firstID, secondID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	var trips []string
	for rows.Next() {
		var tripID string
		if err := rows.Scan(&tripID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		trips = append(trips, tripID)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", rows.Err())
	}
	return trips, nil
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
