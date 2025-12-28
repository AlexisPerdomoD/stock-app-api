package cockroachdb

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {

	dbName := os.Getenv("CR_DB")
	if dbName == "" {
		return nil, fmt.Errorf("CR_DB is empty")
	}

	host := os.Getenv("CR_HOST")
	if host == "" {
		return nil, fmt.Errorf("CR_HOST is empty")
	}

	port := os.Getenv("CR_PORT")
	if port == "" {
		return nil, fmt.Errorf("CR_PORT is empty")
	}

	user := os.Getenv("CR_USER")
	if user == "" {
		return nil, fmt.Errorf("CR_USER is empty")
	}

	password := os.Getenv("CR_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("CR_PASSWORD is empty")
	}

	ssl := os.Getenv("CR_SSL")
	if ssl == "" {
		return nil, fmt.Errorf("CR_SSL is empty")
	}

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, ssl)
	return sqlx.Open("postgres", connection)
}
