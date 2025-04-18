/* All rights and lefts reserved */
package main

import (
	"log"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/controller"
	cockroachdb "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistance/cockroach-db"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/scheduler"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/*
1) Instance db (done)
2) Inject db on repositories implementations (done)
3) Inject repositories and services on usecases (done)
4) Inject usecases on controllers (working)
5) Start server (done)
6) Map controllers routes (working)
7) Pray Golang and start
*/
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err.Error())
	}

	db := cockroachdb.NewDB()
	/* Repositories */
	// mr := cockroachdb.NewMarketRepository(db)
	// cr := cockroachdb.NewCompanyRepository(db)
	// br := cockroachdb.NewBrokerageRepository(db)
	sr := cockroachdb.NewStockRepository(db)
	// rr := cockroachdb.NewRecommendationRepository(db)
	//ur := cockroachdb.NewUserRepository(db)

	/* Usecases */
	getStocksUC := usecase.NewGetStocksUseCase(sr)
	registerStocksUC := usecase.NewRegisterStocksUseCase(sr)

	// getRecommendationByStockUC := usecase.NewGetRecommendationsByStockUseCase(sr, rr)

	// loginUserUC := usecase.NewLoginUseCase(ur)
	// registerUserUC := usecase.NewRegisterUserUseCase(ur)
	// registerUserStockUC := usecase.NewRegisterUserStockUseCase(ur)
	// removeUserStockUC := usecase.NewRemoveUserStockUserCase(ur)

	/* Controllers */
	stockController := controller.NewStockController(getStocksUC)

	/* Start server */
	router := gin.Default()
	/* Set routes */
	stockController.SetRoutes(router)

	if err := router.Run(":3000"); err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("[INFO] Server started")

	/* StockSources */
	mainSSource := service.NewMainSourceStockService(false)

	/* Scheduler jobs */
	scheduler := scheduler.New()
	interval := time.Hour * 24
	timeout := time.Minute * 3
	scheduler.AddStockSourceService(mainSSource, registerStocksUC, timeout, &interval)
	scheduler.StartOnBackground()
}
