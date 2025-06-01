package store

import (
	"context"
	"fmt"
	"os"

	"github.com/Hajdudev/ecoDatabase/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceID struct {
	ServiceID bool `db:"service_id" json:"service_id"`
}

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}

type DatabaseStore interface {
	GetUserByID(id string) (*models.User, error)
	GetRoutesById(firstID []string, secondID []string, ch chan<- map[string]models.TripHash) error
	GetStopInfo(stopID string, ch chan<- models.Stop) error
	GetStopTimesInfo(firstID []string, secondID []string, date string, ch chan<- []models.TempStop) error
	GetStopsID(name string, ch chan<- []string) error
	GetCalendarType(date string, ch chan<- string) error
}

func (pg *PostgresStore) GetCalendarType(date string, ch chan<- string) error {
	query := `SELECT service_id FROM calendar_dates WHERE date = $1`
	var serviceID string
	err := pg.db.QueryRow(context.Background(), query, date).Scan(&serviceID)
	if err != nil {
		ch <- ""
		return err
	}
	ch <- serviceID
	return nil
}

func (pg *PostgresStore) GetStopInfo(stopID string, ch chan<- models.Stop) error {
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

func (pg *PostgresStore) GetStopTimesInfo(firstID []string, secondID []string, date string, ch chan<- []models.TempStop) error {
	query := `
SELECT 
			t1.trip_id,
			t1.stop_id AS from_stop_id,
			t1.departure_time AS from_departure_time,
			t2.stop_id AS to_stop_id,
			t2.departure_time AS to_departure_time
		FROM 
			stop_times t1
		JOIN 
			stop_times t2
			ON t1.trip_id = t2.trip_id
		JOIN 
			trips tr
			ON t1.trip_id = tr.trip_id
		JOIN 
			calendar_dates cd
			ON tr.service_id = cd.service_id
		WHERE 
			t1.stop_id = ANY($1)
			AND t2.stop_id = ANY($2)
			AND cd.date = $3
	`

	firstArray := pgtype.Array[string]{
		Elements: firstID,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(firstID)), LowerBound: 1}},
		Valid:    true,
	}

	secondArray := pgtype.Array[string]{
		Elements: secondID,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(secondID)), LowerBound: 1}},
		Valid:    true,
	}

	rows, err := pg.db.Query(context.Background(), query, &firstArray, &secondArray)
	if err != nil {
		ch <- nil
		return err
	}
	defer rows.Close()

	var trips []models.TempStop
	for rows.Next() {
		var trip models.TempStop
		if err := rows.Scan(&trip.TripID, &trip.FromStopID, &trip.FromDepartureTime, &trip.ToStopID, &trip.ToDepartureTime); err != nil {
			ch <- nil
			return err
		}
		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		ch <- nil
		return err
	}

	ch <- trips
	return nil
}

func (pg *PostgresStore) GetStopsID(name string, ch chan<- []string) error {
	query := `SELECT stop_id FROM stops WHERE stop_name = $1`
	rows, err := pg.db.Query(context.Background(), query, name)
	if err != nil {
		ch <- nil
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			ch <- nil
			return err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		ch <- nil
		return err
	}

	ch <- ids
	return nil
}

func (pg *PostgresStore) GetRoutesById(firstID []string, secondID []string, ch chan<- map[string]models.TripHash) error {
	query := `
    SELECT trip_id, trip_headsign, service_id
    FROM trips
    WHERE trip_id IN (
        SELECT st1.trip_id
        FROM stop_times st1
        JOIN stop_times st2
          ON st1.trip_id = st2.trip_id
         AND st1.stop_id = ANY($1)
         AND st2.stop_id = ANY($2)
         AND st1.stop_sequence < st2.stop_sequence
    )
  `

	firstArray := pgtype.Array[string]{
		Elements: firstID,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(firstID)), LowerBound: 1}},
		Valid:    true,
	}

	secondArray := pgtype.Array[string]{
		Elements: secondID,
		Dims:     []pgtype.ArrayDimension{{Length: int32(len(secondID)), LowerBound: 1}},
		Valid:    true,
	}

	rows, err := pg.db.Query(context.Background(), query, &firstArray, &secondArray)
	if err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		ch <- nil
		return err
	}
	defer rows.Close()

	trips := make(map[string]models.TripHash)
	for rows.Next() {
		var tripID string
		var trip models.TripHash
		if err := rows.Scan(&tripID, &trip.Headsign, &trip.ServiceID); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			ch <- nil
			return err
		}
		trips[tripID] = trip
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error during rows iteration: %v\n", err)
		ch <- nil
		return err
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
