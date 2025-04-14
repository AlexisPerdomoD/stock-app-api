/* All rights and lefts reserved */
package main

import (
	"fmt"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	cockroachdb "github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistance/cockroach-db"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service"
	"github.com/joho/godotenv"
)

/*
1) Instance db (done)
2) Inject db on repositories implementations (done)
3) Inject repositories and services on usecases (done)
4) Inject usecases on controllers
5) Start server
6) Map controllers routes
7) Pray Golang and start
*/
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err.Error())
	}

	db := cockroachdb.NewDB()
	if err := cockroachdb.Migrate(db); err != nil {
		log.Fatalln(err.Error())
	}
	/* Repositories */

	mr := cockroachdb.NewMarketRepository(db)
	cr := cockroachdb.NewCompanyRepository(db)
	br := cockroachdb.NewBrokerageRepository(db)
	sr := cockroachdb.NewStockRepository(db)
	rr := cockroachdb.NewRecommendationRepository(db)
	ur := cockroachdb.NewUserRepository(db)

	/* Usecases */

	getStocksUC := usecase.NewGetStocksUseCase(sr)
	registerStocksUC := usecase.NewRegisterStocksUseCase(sr, cr, mr, br, rr)

	getRecommendationByStockUC := usecase.NewGetRecommendationsByStockUseCase(sr, rr)

	loginUserUC := usecase.NewLoginUseCase(ur)
	registerUserUC := usecase.NewRegisterUserUseCase(ur)
	registerUserStockUC := usecase.NewRegisterUserStockUseCase(ur)
	removeUserStockUC := usecase.NewRemoveUserStockUserCase(ur)

	/* StockSources */

	mainStockSource := service.NewMainSourceStockData()

	fmt.Printf("GetStocksUC: %+v\n", getStocksUC)
	fmt.Printf("RegisterStocksUC: %+v\n", registerStocksUC)
	fmt.Printf("GetRecommendationByStockUC: %+v\n", getRecommendationByStockUC)
	fmt.Printf("LoginUserUC: %+v\n", loginUserUC)
	fmt.Printf("RegisterUserUC: %+v\n", registerUserUC)
	fmt.Printf("RegisterUserStockUC: %+v\n", registerUserStockUC)
	fmt.Printf("RemoveUserStockUC: %+v\n", removeUserStockUC)
	fmt.Printf("MainStockSource: %+v\n", mainStockSource)

}
