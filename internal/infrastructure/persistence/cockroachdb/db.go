package cockroachdb

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("CR_HOST"),
		os.Getenv("CR_PORT"),
		os.Getenv("CR_USER"),
		os.Getenv("CR_PASSWORD"),
		os.Getenv("CR_DB"),
		os.Getenv("CR_SSL"),
	)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil || db == nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// setutp db config
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic DB from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	// sqlDB.SetConnMaxLifetime(0)
	// sqlDB.SetMaxOpenConns(20)
	return db
}

func Migrate(db *gorm.DB) error {
	if db == nil {
		log.Fatalln("bad impl: not db provided to migrate script")
	}

	if os.Getenv("CR_RUN_MIGRATE") != "TRUE" {
		log.Println("AutoMigration was omitted")
		return nil
	}

	return db.AutoMigrate(
		&marketRecord{},
		&companyRecord{},
		&brokerageRecord{},
		&stockRecord{},
		&recommendationRecord{},
		&userRecord{},
	)
}
