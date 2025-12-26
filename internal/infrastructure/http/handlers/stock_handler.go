package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	getStocksUC *usecases.GetStocks
	getStockUC  *usecases.GetStock
}

func NewStockController(
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

func (sc *StockHandler) GetStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID required",
		})
		return
	}

	parsedStockID, err := strconv.Atoi(stockID)
	if err != nil || parsedStockID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID invalid",
		})
		return
	}
	userID := c.GetUint("user_id")
	ctx := c.Request.Context()
	stock, err := sc.getStockUC.Execute(ctx, uint(parsedStockID), &userID)

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

	stocks, err := sc.getStocksUC.Execute(ctx, filters, nil)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (sc *StockHandler) SetRoutes(r *gin.Engine) {
	group := r.Group("/stocks")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("", sc.GetStocksHandler)
	group.GET("/:stockID", sc.GetStockHandler)
}
