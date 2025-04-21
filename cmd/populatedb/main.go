package main

import (
	"context"
	"log"
	"time"

	cockroachdb "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	db := cockroachdb.NewDB()

	if err := cockroachdb.Migrate(db); err != nil {
		log.Fatalln(err.Error())
	}

	mainService := service.NewMainSourceStockService(true)
	sr := cockroachdb.NewStockRepository(db)

	payload, err := mainService.Get(ctx, nil)

	if err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	log.Printf("[populatedb] starting insert of %d stocks", len(payload))

	if err := sr.Register(ctx, payload); err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	log.Printf("[populatedb] done with %d stocks", len(payload))
}
