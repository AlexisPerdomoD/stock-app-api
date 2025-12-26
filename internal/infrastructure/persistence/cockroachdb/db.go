package cockroachdb

import (
	"database/sql"
	"fmt"
	"os"
)

func NewDB() *sql.DB {
	_ = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("CR_HOST"),
		os.Getenv("CR_PORT"),
		os.Getenv("CR_USER"),
		os.Getenv("CR_PASSWORD"),
		os.Getenv("CR_DB"),
		os.Getenv("CR_SSL"),
	)

	return nil
}

func Migrate() error {

	return nil
}
