package controller

import (
	"log"
	"net/http"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/dto"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

type StockController struct {
	getStocksUseCase *usecase.GetStocksUseCase
}

func (sc *StockController) GetStocksHandler(c *gin.Context) {
	ctx := c.Request.Context()
	filters := dto.MapGetStocksFilter(c)
	
	stocks, err := sc.getStocksUseCase.Execute(ctx, *filters, nil)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (sc *StockController) SetRoutes(r *gin.Engine) {
	group := r.Group("/stocks")

	group.GET("", sc.GetStocksHandler)
}

func NewStockController(getStocksUseCase *usecase.GetStocksUseCase) *StockController {
	if getStocksUseCase == nil {
		log.Fatalln("bad impl: GetStocksUseCase was nil for NewStockController")
	}

	return &StockController{getStocksUseCase}
}
