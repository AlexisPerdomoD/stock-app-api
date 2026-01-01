/* All rights and lefts reserved */
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/handlers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/scheduler"
	servicesimpl "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/*
1) Instance db
2) Inject db on repositories implementations
3) Inject repositories and services on usecases
4) Inject usecases on controllers
5) Map controllers routes
6) Set Cron jobs (working)
7) Start server
*/
func main() {
	// ENV VARS
	_ = godotenv.Load()

	// DATABASE
	db, err := cockroachdb.NewDB()
	if err != nil {
		panic(fmt.Sprintf("Error creating db due to %v", err))
	}

	// SERVICES
	mainSlogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	unitOfWorkFactory := servicesimpl.NewSqlxUnitOfWorkFactory(db)

	mainDataSourceService := servicesimpl.NewMainSourceStockService(true)
	cnnDataSourceService := servicesimpl.NewCnnStockSourceService()
	const mainsourcekey, cnnsourcekey = "main", "cnn"
	dataSources := map[string]services.DataSourceService{
		mainsourcekey: mainDataSourceService,
		cnnsourcekey:  cnnDataSourceService,
	}

	// REPOSITORIES

	userRepository := cockroachdb.NewUserRepository(db)
	marketRepository := cockroachdb.NewMarketRepository(db)
	companyRepository := cockroachdb.NewCompanyRepository(db)
	stockRepository := cockroachdb.NewStockRepository(db)
	stockRegisterRepository := cockroachdb.NewStockRegisterRepository(db)
	recommendationRepository := cockroachdb.NewRecommendationRepository(db)

	// USE CASES
	getStocksUC := usecases.NewGetStocks(stockRepository)
	getStockUC := usecases.NewGetStock(
		userRepository,
		marketRepository,
		companyRepository,
		stockRepository,
		stockRegisterRepository,
	)
	registerStocksUC := usecases.NewRegisterStocks(unitOfWorkFactory, dataSources)
	getRecommendationByStockUC := usecases.NewGetRecommendationsByStock(stockRepository, recommendationRepository)
	loginUserUC := usecases.NewLogin(userRepository, mainSlogger)
	registerUserUC := usecases.NewRegisterUser(userRepository, mainSlogger)
	registerUserStockUC := usecases.NewRegisterUserStock(userRepository)
	removeUserStockUC := usecases.NewRemoveUserStock(userRepository)

	// HANDLERS
	stockHandler := handlers.NewStockHandler(getStocksUC, getStockUC)
	recommendationHandler := handlers.NewRecommendationHandler(getRecommendationByStockUC)
	userHandler := handlers.NewUserHandler(
		getStocksUC,
		registerUserUC,
		loginUserUC,
		registerUserStockUC,
		removeUserStockUC,
	)

	// ROUTES
	r := gin.Default()
	// TODO: Implement cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	r.POST("/api/v1/login", userHandler.LoginUserHandler)

	userGroup := r.Group("/api/v1/users")
	userGroup.Use(middleware.UserSessionMiddleware)
	userGroup.POST("", userHandler.RegisterUserHandler)
	userGroup.GET("/stocks", userHandler.GetStocksHandler)
	userGroup.POST("/stocks/:stockID", userHandler.RegisterStockHandler)
	userGroup.DELETE("/stocks/:stockID", userHandler.RemoveStockHandler)

	stockGroup := r.Group("/api/v1/stocks")
	stockGroup.Use(middleware.UserSessionMiddleware)
	stockGroup.GET("", stockHandler.GetStocksHandler)
	stockGroup.GET("/:stockID", stockHandler.GetStockHandler)

	recommendationGroup := r.Group("/api/v1/recommendations")
	recommendationGroup.Use(middleware.UserSessionMiddleware)
	recommendationGroup.GET("/:stockID", recommendationHandler.GetRecommendationsByStockHandler)

	// SCHEDULER
	scheduler := scheduler.New()
	interval := time.Minute * 10
	timeout := time.Minute * 3
	scheduler.AddStockSourceService(
		mainsourcekey,
		registerStocksUC,
		timeout,
		&interval,
	)
	scheduler.AddStockSourceService(
		cnnsourcekey,
		registerStocksUC,
		timeout,
		&interval,
	)
	scheduler.StartOnBackground()

	PORT := fmt.Sprintf(":%v", os.Getenv("SERVER_PORT"))
	if err := r.Run(PORT); err != nil {
		slog.Error("Error running server", "err", err)
		os.Exit(1)
	}

}
