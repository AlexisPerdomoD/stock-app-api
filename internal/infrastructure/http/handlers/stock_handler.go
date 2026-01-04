package handlers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StockHandler struct {
	getStocks         *usecases.GetStocks
	getStock          *usecases.GetStock
	registerUserStock *usecases.RegisterUserStock
	removeUserStock   *usecases.RemoveUserStock
}

// GetStockHandler godoc
// @Summary Obtener stock
// @Description Obtener un stock por su ID
// @Tags stocks
// @Produce json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param stockID path string true "ID del stock"
// @Success 200 {object} models.PopulatedStockView
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/stocks/{stockID} [get]
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

// GetStocksHandler godoc
// @Summary Obtener stocks paginados
// @Description Retorna una lista paginada en base a los filtros aplicados
// @Tags stocks
// @Produce json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param search query string false "Búsqueda por ticker o nombre"
// @Param orderby query string false "Ordenar resultados. Valores: tendency-asc, tendency-desc, price-asc, price-desc, ticker-asc, ticker-desc, date"
// @Param greater query number false "Precio mínimo"
// @Param lower query number false "Precio máximo"
// @Param market query number false "ID de mercado fuente de datos"
// @Param page query int false "Número de página. Default: 1"
// @Param size query int false "Tamaño de página. Default: 20"
// @Success 200 {object} models.StockViewPaginated
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/stocks [get]
func (sc *StockHandler) GetStocksHandler(c *gin.Context) {
	ctx := c.Request.Context()
	filters := mappers.MapGetStocksFilter(c)

	stocks, err := sc.getStocks.Execute(ctx, filters, nil)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// GetStocksHandler godoc
// @Summary Obtener stocks favoritos del usuario
// @Description Retorna la lista de stocks asociados al usuario autenticado
// @Tags stocks
// @Produce json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param search query string false "Búsqueda por ticker o nombre"
// @Param orderby query string false "Ordenar resultados. Valores: tendency-asc, tendency-desc, price-asc, price-desc, ticker-asc, ticker-desc, date"
// @Param greater query number false "Precio mínimo"
// @Param lower query number false "Precio máximo"
// @Param market query number false "ID de mercado fuente de datos"
// @Param page query int false "Número de página. Default: 1"
// @Param size query int false "Tamaño de página. Default: 20"
// @in header
// @name authorization
// @Success 200 {object} models.StockViewPaginated
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/stocks/favorites [get]
func (sc *StockHandler) GetStocksByUserHandler(c *gin.Context) {
	userID := c.GetUint64("user_id")
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := c.Request.Context()
	filters := mappers.MapGetStocksFilter(c)

	stocks, err := sc.getStocks.Execute(ctx, filters, &userID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// RegisterStockHandler godoc
// @Summary Registrar stock al usuario
// @Description Asocia un stock existente al usuario autenticado
// @Tags stocks
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @in header
// @name authorization
// @Param stockID path int true "ID del stock"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 404 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/stocks/{stockID}/favorites [post]
func (uc *StockHandler) RegisterStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	parseStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	userID := c.GetUint64("user_id")

	ctx := c.Request.Context()
	if err := uc.registerUserStock.Execute(ctx, userID, parseStockID); err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": true, "message": "user stock registered"})
}

// RemoveStockHandler godoc
// @Summary Eliminar stock del usuario
// @Description Remueve un stock asociado al usuario autenticado
// @Tags stocks
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @in header
// @name authorization
// @Param stockID path int true "ID del stock"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 404 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/stocks/{stockID}/favorites [delete]
func (uc *StockHandler) RemoveStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	parseStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	userID := c.GetUint64("user_id")

	ctx := c.Request.Context()
	if err := uc.removeUserStock.Execute(ctx, userID, parseStockID); err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "user stock removed"})
}

func NewStockHandler(
	getStocksUC *usecases.GetStocks,
	getStockUC *usecases.GetStock,
	registerUserStockUC *usecases.RegisterUserStock,
	removeUserStockUC *usecases.RemoveUserStock,
) *StockHandler {
	if getStocksUC == nil || getStockUC == nil || registerUserStockUC == nil || removeUserStockUC == nil {
		panic("[StockController]: one or more usecases provided as nil")
	}

	return &StockHandler{getStocksUC, getStockUC, registerUserStockUC, removeUserStockUC}
}
