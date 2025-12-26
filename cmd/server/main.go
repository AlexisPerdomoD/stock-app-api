/* All rights and lefts reserved */
package main

import (
	"fmt"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/handlers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/scheduler"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
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
	_ = godotenv.Load()

	db := cockroachdb.NewDB()

	sr := cockroachdb.NewStockRepository(db)
	rr := cockroachdb.NewRecommendationRepository(db)
	ur := cockroachdb.NewUserRepository(db)

	getStocksUC := usecases.NewGetStocks(sr)
	getStockUC := usecases.NewGetStock(sr)
	registerStocksUC := usecases.NewRegisterStocks(sr)
	getRecommendationByStockUC := usecases.NewGetRecommendationsByStock(sr, rr)
	loginUserUC := usecases.NewLogin(ur)
	registerUserUC := usecases.NewRegisterUser(ur)
	registerUserStockUC := usecases.NewRegisterUserStock(ur)
	removeUserStockUC := usecases.NewRemoveUserStock(ur)

	stockController := handlers.NewStockHandler(getStocksUC, getStockUC)
	recommendationController := handlers.NewRecommendationHandler(getRecommendationByStockUC)
	userController := handlers.NewUserHandler(
		getStocksUC,
		registerUserUC,
		loginUserUC,
		registerUserStockUC,
		removeUserStockUC,
	)

	router := gin.Default()
	// TODO: Implement cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	stockController.SetRoutes(router)
	recommendationController.SetRoutes(router)
	userController.SetRoutes(router)

	scheduler := scheduler.New()

	mainSSource := service.NewMainSourceStockService(false)
	interval := time.Hour * 24
	timeout := time.Minute * 3
	scheduler.AddStockSourceService(
		mainSSource,
		registerStocksUC,
		timeout,
		&interval,
	)
	scheduler.StartOnBackground()

	PORT := fmt.Sprintf(":%v", os.Getenv("SERVER_PORT"))
	if err := router.Run(PORT); err != nil {
		log.Fatalln(err.Error())
	}

}
