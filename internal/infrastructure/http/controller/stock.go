package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/dto"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

type StockController struct {
	getStocksUC *usecase.GetStocksUseCase
	getStockUC  *usecase.GetStockUseCase
}

func NewStockController(
	getStocksUC *usecase.GetStocksUseCase,
	getStockUC *usecase.GetStockUseCase,
) *StockController {
	if getStocksUC == nil {
		log.Fatalln("[StockController]: getStocksUC provided as nil")
	}

	if getStockUC == nil {
		log.Fatalln("[StockController]: getStockUC provided as nil")
	}

	return &StockController{getStocksUC, getStockUC}
}

func (sc *StockController) GetStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID invalid",
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
		res := pkg.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (sc *StockController) GetStocksHandler(c *gin.Context) {
	ctx := c.Request.Context()
	filters := dto.MapGetStocksFilter(c)

	stocks, err := sc.getStocksUC.Execute(ctx, filters, nil)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (sc *StockController) SetRoutes(r *gin.Engine) {
	group := r.Group("/stocks")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("", sc.GetStocksHandler)
	group.GET("/:stockID", sc.GetStockHandler)
}
