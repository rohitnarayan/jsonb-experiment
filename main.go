package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

type Order struct {
	Total              int       `json:"total"`
	LastOrderTimestamp time.Time `json:"lastOrderTimestamp"`
}

type Orders struct {
	Food      Order `json:"food,omitempty"`
	Transport Order `json:"transport,omitempty"`
}

type Metadata struct {
	Orders Orders `json:"orders,omitempty"`
	Locale string `json:"locale,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Metadata string `json:"metadata"`
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

	metadata := Metadata{
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
		Metadata: string(metadataJSON),
	}

	_, err = db.Exec("INSERT INTO jsonb_experiment (id, name, phone, email, metadata) VALUES ($1, $2, $3, $4, $5)", userToInsert.ID, userToInsert.Name, userToInsert.Phone, userToInsert.Email, userToInsert.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	// Example of reading data from the database.
	var retrievedUser User

	if err := db.Get(&retrievedUser, "SELECT id, name, phone, email, metadata FROM jsonb_experiment WHERE id=$1", id); err != nil {
		log.Fatal(err)
	}

	var jsonMeta Metadata
	if err := json.Unmarshal([]byte(retrievedUser.Metadata), &jsonMeta); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved User: %+v\n\n", retrievedUser)
	fmt.Printf("User metadata: %+v\n\n", jsonMeta)
}
