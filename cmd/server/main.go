/* All rights and lefts reserved */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/handlers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file due to %v", err)
	}

	// REPOSITORIES
	db := cockroachdb.NewDB()
	stockRepository := cockroachdb.NewStockRepository(db)
	recommendationRepository := cockroachdb.NewRecommendationRepository(db)
	userRepository := cockroachdb.NewUserRepository(db)

	// USE CASES
	getStocksUC := usecases.NewGetStocks(stockRepository)
	getStockUC := usecases.NewGetStock(stockRepository)
	// registerStocksUC := usecases.NewRegisterStocks(sr)
	getRecommendationByStockUC := usecases.NewGetRecommendationsByStock(stockRepository, recommendationRepository)
	loginUserUC := usecases.NewLogin(userRepository)
	registerUserUC := usecases.NewRegisterUser(userRepository)
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

	// scheduler := scheduler.New()
	//
	// mainSSource := service.NewMainSourceStockService(false)
	// interval := time.Hour * 24
	// timeout := time.Minute * 3
	// scheduler.AddStockSourceService(
	// 	mainSSource,
	// 	registerStocksUC,
	// 	timeout,
	// 	&interval,
	// )
	// scheduler.StartOnBackground()
	//
	PORT := fmt.Sprintf(":%v", os.Getenv("SERVER_PORT"))
	if err := r.Run(PORT); err != nil {
		log.Fatalln(err.Error())
	}

}
