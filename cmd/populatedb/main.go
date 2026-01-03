package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	servicesimpl "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/services"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()
	servicename := flag.Arg(0)
	const mainServiceName, cnnServiceName = "main", "cnn"

	switch servicename {
	case mainServiceName:
	case cnnServiceName:
	default:
		panic(fmt.Sprintf("invalid service name %s", servicename))
	}

	_ = godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	db, err := cockroachdb.NewDB()
	if err != nil {
		slog.Log(ctx, slog.LevelError, "error connecting db", "err", err)
		os.Exit(1)
	}

	httpServiceClient := http.Client{Timeout: time.Second * 10}
	mainService := servicesimpl.NewMainSourceStockService(&httpServiceClient, true)
	cnnService := servicesimpl.NewCnnStockSourceService(&httpServiceClient)

	serviceProvider := make(map[string]services.DataSourceService)
	serviceProvider[mainServiceName] = mainService
	serviceProvider[cnnServiceName] = cnnService

	unitOfWorkFactory := cockroachdb.NewUnitOfWorkFactory(db, slog.Default().With("cockroachdb", "UnitOfWorkFactory"))

	registerStocks := usecases.NewRegisterStocks(unitOfWorkFactory, serviceProvider, slog.Default().With("usecase", "RegisterStocks"))

	count, err := registerStocks.Execute(ctx, servicename, nil)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "error resolving registers", "err", err)
		os.Exit(1)
	}

	slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("[populatedb] done with %d stocks", count))
}
