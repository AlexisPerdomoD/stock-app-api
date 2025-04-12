package cockroachdb

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(conn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(conn))

	if err != nil || db == nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// setutp db config
	// maybe in another func with *gorm.DB as param (???) :b
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic DB from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	// sqlDB.SetConnMaxLifetime(0)
	// sqlDB.SetMaxOpenConns(20)
	return db
}
