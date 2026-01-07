package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/handlers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/scheduler"
	servicesimpl "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/services"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/services/mock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		panic(fmt.Sprintf("Error creating db due to %+v", err))
	}

	if err = cockroachdb.MigrateUp(db.DB); err != nil {
		panic(fmt.Sprintf("Error migrating db up due to %+v", err))
	}

	// SERVICES
	httpServiceClient := http.Client{Timeout: time.Second * 10}

	mainSlogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mainDataSourceService := servicesimpl.NewMainSourceStockService(&httpServiceClient, mainSlogger.With("main", "source"))
	cnnDataSourceService := servicesimpl.NewCnnStockSourceService(&httpServiceClient, mainSlogger.With("cnn", "source"))
	mockDataSourceService := mock.NewMockSourceStockService()
	const mainsourcekey, cnnsourcekey, mocksourcekey string = "principal", "cnn", "mock"
	dataSources := map[string]services.DataSourceService{
		mainsourcekey: mainDataSourceService,
		cnnsourcekey:  cnnDataSourceService,
		mocksourcekey: mockDataSourceService,
	}

	// REPOSITORIES

	userRepository := cockroachdb.NewUserRepository(db)
	marketRepository := cockroachdb.NewMarketRepository(db)
	companyRepository := cockroachdb.NewCompanyRepository(db)
	stockRepository := cockroachdb.NewStockRepository(db)
	stockRegisterRepository := cockroachdb.NewStockRegisterRepository(db)
	recommendationRepository := cockroachdb.NewRecommendationRepository(db)
	stockStatsRepository := cockroachdb.NewStockTendencyStatRepository(db)
	unitOfWorkFactory := cockroachdb.NewUnitOfWorkFactory(db, mainSlogger.With("cockroachdb", "UnitOfWorkFactory"))
	// USE CASES
	getMarketsUC := usecases.NewGetMarkets(marketRepository)
	getStocksUC := usecases.NewGetStocks(stockRepository)
	getStockUC := usecases.NewGetStock(
		userRepository,
		marketRepository,
		companyRepository,
		stockRepository,
		stockRegisterRepository,
	)
	registerStocksUC := usecases.NewRegisterStocks(unitOfWorkFactory, dataSources, mainSlogger.With("usecase", "RegisterStocks"))
	getRecommendationByStockUC := usecases.NewGetRecommendations(recommendationRepository)
	loginUserUC := usecases.NewLogin(userRepository, mainSlogger)
	registerUserUC := usecases.NewRegisterUser(userRepository, mainSlogger)
	registerUserStockUC := usecases.NewRegisterUserStock(userRepository)
	removeUserStockUC := usecases.NewRemoveUserStock(userRepository)
	threeMonths := time.Hour * 24 * 30 * 3
	getStockRegistersByStockDateRangedUC := usecases.NewGetStockRegistersByStockDateRanged(stockRegisterRepository, threeMonths)
	getLastStockRegistersByStockUC := usecases.NewGetLastStockRegistersByStock(stockRegisterRepository)
	getStockTendencyStatByStockUC := usecases.NewGetStockRegistersStatsByStock(stockRepository, stockStatsRepository)

	// HANDLERS
	marketHandler := handlers.NewMarketHandler(getMarketsUC)
	stockHandler := handlers.NewStockHandler(getStocksUC, getStockUC, registerUserStockUC, removeUserStockUC)
	stockRegisterHandler := handlers.NewStockRegisterHandler(getStockRegistersByStockDateRangedUC, getLastStockRegistersByStockUC, getStockTendencyStatByStockUC)
	recommendationHandler := handlers.NewRecommendationHandler(getRecommendationByStockUC)
	userHandler := handlers.NewUserHandler(registerUserUC, loginUserUC)

	// ROUTES
	r := gin.Default()
	// TODO: Implement cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/login", userHandler.LoginUserHandler)

	userGroup := r.Group("/api/v1/users")
	userGroup.POST("", userHandler.RegisterUserHandler)

	marketGroup := r.Group("/api/v1/markets")
	marketGroup.Use(middleware.UserSessionMiddleware)
	marketGroup.GET("", marketHandler.GetMarketsHandler)

	stockGroup := r.Group("/api/v1/stocks")
	stockGroup.Use(middleware.UserSessionMiddleware)
	stockGroup.GET("", stockHandler.GetStocksHandler)
	stockGroup.GET("/favorites", stockHandler.GetStocksByUserHandler)
	stockGroup.POST("/favorites/:stockID", stockHandler.RegisterStockHandler)
	stockGroup.DELETE("/favorites/:stockID", stockHandler.RemoveStockHandler)
	stockGroup.GET("/:stockID", stockHandler.GetStockHandler)

	stockRegisterGroup := stockGroup.Group("/:stockID/registers")
	stockRegisterGroup.GET("", stockRegisterHandler.GetStockRegistersByStockDateRangedHandler)
	stockRegisterGroup.GET("/last", stockRegisterHandler.GetLastStockRegistersByStockHandler)
	stockRegisterGroup.GET("/tendency", stockRegisterHandler.GetStockTendencyStatByStockHandler)

	recommendationGroup := r.Group("/api/v1/recommendations")
	recommendationGroup.Use(middleware.UserSessionMiddleware)
	recommendationGroup.GET("/:stockID", recommendationHandler.GetRecommendationsByStockHandler)

	// SCHEDULER
	scheduler := scheduler.New()
	mainInterval := time.Hour * 24
	cnnInterval := time.Hour
	mockInterval := time.Minute * 45
	timeout := time.Minute * 3
	scheduler.AddStockSourceService(
		mainsourcekey,
		registerStocksUC,
		timeout,
		&mainInterval,
	)

	scheduler.AddStockSourceService(
		cnnsourcekey,
		registerStocksUC,
		timeout,
		&cnnInterval,
	)
	scheduler.AddStockSourceService(
		mocksourcekey,
		registerStocksUC,
		timeout,
		&mockInterval,
	)
	scheduler.StartOnBackground()

	PORT := fmt.Sprintf(":%v", os.Getenv("SERVER_PORT"))
	if err := r.Run(PORT); err != nil {
		slog.Error("Error running server", "err", err)
		os.Exit(1)
	}

}
