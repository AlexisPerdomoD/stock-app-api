package handlers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type StockHandler struct {
	getStocks *usecases.GetStocks
	getStock  *usecases.GetStock
}

func (sc *StockHandler) GetStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {

		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	parsedStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil || parsedStockID <= 0 {
		res := mappers.MapHttpErr(pkg.BadRequest("invalid stockID provided"))
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	userID := c.GetUint64("user_id")
	ctx := c.Request.Context()
	stock, err := sc.getStock.Execute(ctx, parsedStockID, &userID)

	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (sc *StockHandler) GetStocksHandler(c *gin.Context) {
	ctx := c.Request.Context()
	filters := mappers.MapGetStocksFilter(c)

	res, err := sc.getStocks.Execute(ctx, filters, nil)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (sc *StockHandler) SetRoutes(r *gin.Engine) {
	group := r.Group("/stocks")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("", sc.GetStocksHandler)
	group.GET("/:stockID", sc.GetStockHandler)
}

func NewStockHandler(
	getStocksUC *usecases.GetStocks,
	getStockUC *usecases.GetStock,
) *StockHandler {
	if getStocksUC == nil {
		log.Fatalln("[StockController]: getStocksUC provided as nil")
	}

	if getStockUC == nil {
		log.Fatalln("[StockController]: getStockUC provided as nil")
	}

	return &StockHandler{getStocksUC, getStockUC}
}
