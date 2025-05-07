package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Hajdudev/ecoDatabase/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

type DatabaseStore interface {
	GetRoutes()
}

func (pg *PostgresStore) GetRoutes() {
	query := "SELECT id, created_at, email, name, image, recent_rided FROM users"

	rows, err := pg.db.Query(query)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Users found in database:")

	for rows.Next() {
		var user models.User
		var recentRidedBytes []byte

		err := rows.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Name, &user.Image, &recentRidedBytes)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		if len(recentRidedBytes) > 0 {
			json.Unmarshal(recentRidedBytes, &user.RecentRided)
		}

		fmt.Printf("User: %+v\n", user)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
		return
	}

	fmt.Println("Query completed successfully")
}
