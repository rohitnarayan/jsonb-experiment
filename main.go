package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

type Metadata map[string]interface{}

type User struct {
	ID       int      `json:"id" db:"id"`
	Name     string   `json:"name" db:"name"`
	Phone    string   `json:"phone" db:"phone"`
	Email    string   `json:"email" db:"email"`
	Metadata Metadata `json:"metadata" db:"metadata"`
}

// Value method for MetadataJSON
func (m *Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan method for MetadataJSON
func (m *Metadata) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}

func main() {
	// Create a connection string.
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Initialize a database connection using sqlx.
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	/**
	  Orders: Orders{
	  		Food: Order{
	  			Total:              50,
	  			LastOrderTimestamp: time.Now(),
	  		},
	  		Transport: Order{
	  			Total:              30,
	  			LastOrderTimestamp: time.Now(),
	  		},
	  	},
	  	Locale: "en_US",
	*/

	metadata := map[string]interface{}{
		"orders": map[string]interface{}{
			"food": map[string]interface{}{
				"total":              50,
				"lastOrderTimestamp": time.Now().UTC(),
			},
			"transport": map[string]interface{}{
				"total":              150,
				"lastOrderTimestamp": time.Now().UTC(),
			},
		},
		"locale": "en_ID",
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Fatal(err)
	}

	// Example of writing data to the database.
	id := rand.Intn(100000)
	userToInsert := User{
		ID:       id,
		Name:     "RJ",
		Phone:    "+911234567890",
		Email:    "rj@gmail.com",
		Metadata: metadata,
	}

	_, err = db.Exec("INSERT INTO jsonb_experiment (id, name, phone, email, metadata) VALUES ($1, $2, $3, $4, $5)",
		userToInsert.ID, userToInsert.Name, userToInsert.Phone, userToInsert.Email, metadataJSON)
	if err != nil {
		log.Fatal(err)
	}

	// Example of reading data from the database.
	var retrievedUser User
	if err := db.Get(&retrievedUser, "SELECT id, name, phone, email, metadata FROM jsonb_experiment WHERE id=$1", id); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved User: %+v\n", retrievedUser)
}
