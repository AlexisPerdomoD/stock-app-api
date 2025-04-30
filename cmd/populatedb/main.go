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
	cnnService := service.NewCnnStockSourceService()

	sr := cockroachdb.NewStockRepository(db)

	payloadcnn, err := cnnService.Get(ctx, nil)
	if err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	if err := sr.Register(ctx, payloadcnn); err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	payloadmain, err := mainService.Get(ctx, nil)

	if err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	log.Printf("[populatedb] starting insert of %d stocks", len(payloadmain))

	if err := sr.Register(ctx, payloadmain); err != nil {
		log.Fatalf("[populatedb] error: %v", err)
	}

	log.Printf("[populatedb] done with %d stocks", len(payloadmain)+len(payloadcnn))
}
