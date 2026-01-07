package cockroachdb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	migratecockroachdb "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	ssl := os.Getenv("CR_SSL")
	if ssl == "" {
		return nil, fmt.Errorf("CR_SSL is empty")
	}

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		ssl)

	return sqlx.Open("postgres", connection)
}

func MigrateUp(db *sql.DB) error {
	driver, err := migratecockroachdb.WithInstance(db, &migratecockroachdb.Config{})
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := "file://" + filepath.Join(
		wd,
		"internal/infrastructure/db/migrations/cockroachdb",
	)

	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
