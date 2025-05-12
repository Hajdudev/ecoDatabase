package store

import (
	"context"
	"fmt"
	"os"
	"sync"

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
	GetRoutesById(firstID string, secondID string, ch chan<- []models.Trip, wg *sync.WaitGroup) error
	GetStopInfo(stopID string, ch chan<- models.Stop, wg *sync.WaitGroup) error
}

func (pg *PostgresStore) GetStopInfo(stopID string, ch chan<- models.Stop, wg *sync.WaitGroup) error {
	query := `
		SELECT stop_id, stop_code, stop_name, stop_desc, stop_lat, stop_lon
		FROM stops
		WHERE stop_id = $1
	`
	var stop models.Stop

	defer wg.Done()
	err := pg.db.QueryRow(context.Background(), query, stopID).Scan(
		&stop.StopID,
		&stop.StopCode,
		&stop.StopName,
		&stop.StopDesc,
		&stop.StopLat,
		&stop.StopLon,
	)
	if err != nil {
		ch <- models.Stop{}
		return err
	}

	ch <- stop
	return nil
}

func (pg *PostgresStore) GetRoutesById(firstID string, secondID string, ch chan<- []models.Trip, wg *sync.WaitGroup) error {
	query := `
		SELECT trip_headsign, trip_id
		FROM trips 
		WHERE trip_id IN (
			SELECT trip_id
			FROM stop_times
			WHERE stop_id IN ($1, $2)
			GROUP BY trip_id
			HAVING COUNT(DISTINCT stop_id) = 2
		)
	`
	defer wg.Done()

	rows, err := pg.db.Query(context.Background(), query, firstID, secondID)
	if err != nil {
		fmt.Println("Error querying database:", err)
		ch <- nil
		return err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var trip models.Trip
		if err := rows.Scan(&trip.TripHeadsign, &trip.TripID); err != nil {
			fmt.Println("Error scanning row:", err)
			ch <- nil
			return err
		}
		trips = append(trips, trip)
	}

	if rows.Err() != nil {
		fmt.Println("Error during rows iteration:", rows.Err())
		ch <- nil
		return rows.Err()
	}

	ch <- trips
	return nil
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
