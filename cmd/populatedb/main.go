package main

import (
	"context"
	"log"
	"time"

	cockroachdb "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistance/cockroach-db"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	db := cockroachdb.NewDB()

	mainService := service.NewMainSourceStockService(true)
	sr := cockroachdb.NewStockRepository(db)

	payload, err := mainService.Get(ctx, nil)

	if err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	if err := sr.Register(ctx, payload); err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	log.Printf("[populatedb] done with %d stocks", len(payload))

}
