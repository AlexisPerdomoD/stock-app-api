package controller

import (
	"log"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

type StockController struct {
	getStocksUseCase *usecase.GetStocksUseCase
}

func (sc *StockController) GetStocksHandler(c *gin.Context) {
	ctx := c.Request.Context()
	// filter by lower than
	// filter by greater than
	// filter by company name
	// filter by ticker
	// filter by user id
	filters := pkg.PaginationFilter{
		SortBy: map[string]pkg.SortOrder{
			"tendency": pkg.SortOrderAsc,
			"price":    pkg.SortOrderDesc,
		},

		PaginationPage: pkg.PaginationPage{
			Size: 20,
			Page: 1,
		},
		FilterBy: []pkg.FilterByItem{},
	}

	page := c.DefaultQuery("page", "1")
	parsePage, err := strconv.Atoi(page)
	if err == nil {
		filters.PaginationPage.Page = parsePage
	}

	stocks, err := sc.getStocksUseCase.Execute(ctx, filters)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(200, stocks)

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
