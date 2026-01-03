package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
)

type StockRegisterHandler struct {
	getStockRegistersByStockDateRanged *usecases.GetStockRegistersByStockDateRanged

	getLastStockRegistersByStock *usecases.GetLastStockRegistersByStock

	getStockTendencyStatByStock *usecases.GetStockRegistersStatsByStock
}

// GetStockRegistersByStockDateRangedHandler godoc
//
// @Summary      Get stock registers by date range
// @Description  Returns stock register history for a given stock order by date DESC.
// @Description  If `from` and `to` are not provided, defaults to last 30 days.
// @Tags         stock-registers
// @Accept       json
// @Produce      json
// @Security 	 BearerAuth
// @Param 		 authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param        stockID path     int     true  "Stock ID"
// @Param        from    query    string  false "From date (RFC3339)" example(2025-01-01T00:00:00Z)
// @Param        to      query    string  false "To date (RFC3339)"   example(2025-02-01T00:00:00Z)
// @Success      200 {array}  models.StockRegisterView
// @Failure      400 {object} mappers.HttpErrResponse
// @Failure      404 {object} mappers.HttpErrResponse
// @Failure      500 {object} mappers.HttpErrResponse
// @Router       /api/v1/stocks/{stockID}/registers [get]
func (h *StockRegisterHandler) GetStockRegistersByStockDateRangedHandler(c *gin.Context) {
	rawStockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("missing stockID"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	stockID, err := strconv.ParseUint(rawStockID, 10, 64)
	if err != nil || stockID == 0 {
		res := mappers.MapHttpErr(pkg.BadRequest("invalid stockID"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	fromRaw, toRaw := c.Query("from"), c.Query("to")
	var from, to time.Time

	if fromRaw != "" && toRaw != "" {
		from, err = time.Parse(time.RFC3339, fromRaw)
		if err != nil {
			res := mappers.MapHttpErr(pkg.BadRequest(
				fmt.Sprintf("invalid from date: %s, is expected in RFC3339 format (e.g. 2006-01-02T15:04:05Z07:00)", fromRaw)))
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		to, err = time.Parse(time.RFC3339, toRaw)
		if err != nil {
			res := mappers.MapHttpErr(pkg.BadRequest(
				fmt.Sprintf("invalid to date: %s, is expected in RFC3339 format (e.g. 2006-01-02T15:04:05Z07:00)", toRaw)))
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		from = from.UTC()
		to = to.UTC()
	} else if fromRaw != "" || toRaw != "" {
		res := mappers.MapHttpErr(pkg.BadRequest("from and to must be provided together"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	} else {
		from, to = time.Now().AddDate(0, -1, 0).UTC(), time.Now().UTC()
	}

	res, err := h.getStockRegistersByStockDateRanged.Execute(c.Request.Context(), stockID, from, to)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetLastStockRegistersByStockHandler godoc
// @Summary      Get last stock registers ordered by date DESC
// @Description  Returns the last N stock registers for a given stock ordered by date DESC.
// @Description  If `limit` is not provided, defaults to 100.
// @Tags         stock-registers
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param        stockID  path     int   true  "Stock ID"
// @Param        limit    query    int   false "Max number of registers to return" minimum(1) maximum(65535) default(100)
// @Success      200 {array}  models.StockRegisterView
// @Failure      400 {object} mappers.HttpErrResponse
// @Failure      404 {object} mappers.HttpErrResponse
// @Failure      500 {object} mappers.HttpErrResponse
// @Router       /api/v1/stocks/{stockID}/registers/last [get]
func (h *StockRegisterHandler) GetLastStockRegistersByStockHandler(c *gin.Context) {
	rawStockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	stockID, err := strconv.ParseUint(rawStockID, 10, 64)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	var limit uint16
	rawLimit, ok := c.GetQuery("limit")
	if ok {
		limit64, err := strconv.ParseUint(rawLimit, 10, 16)
		if err != nil {
			res := mappers.MapHttpErr(pkg.BadRequest(fmt.Sprintf("limit is invalid: %s, must be a positive int with max value of %d", rawLimit, uint16(math.MaxUint16))))
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}
		limit = uint16(limit64)
	} else {
		limit = 100
	}

	ctx := c.Request.Context()
	res, err := h.getLastStockRegistersByStock.Execute(ctx, stockID, limit)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetStockTendencyStatByStockHandler godoc
// @Summary      Get stock tendency stat by stock
// @Description  Returns the stock tendency stat for a given stock.
// @Tags         stock-registers
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param        stockID  path     int   true  "Stock ID"
// @Success      200 {object}  models.StockTendencyStatView
// @Failure      400 {object} mappers.HttpErrResponse
// @Failure      404 {object} mappers.HttpErrResponse
// @Failure      500 {object} mappers.HttpErrResponse
// @Router       /api/v1/stocks/{stockID}/registers/tendency [get]
func (h *StockRegisterHandler) GetStockTendencyStatByStockHandler(c *gin.Context) {
	rawStockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	stockID, err := strconv.ParseUint(rawStockID, 10, 64)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	ctx := c.Request.Context()
	res, err := h.getStockTendencyStatByStock.Execute(ctx, stockID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func NewStockRegisterHandler(
	getStockRegistersByStockDateRanged *usecases.GetStockRegistersByStockDateRanged,
	getLastStockRegistersByStock *usecases.GetLastStockRegistersByStock,
	getStockTendencyStatByStock *usecases.GetStockRegistersStatsByStock,
) *StockRegisterHandler {
	if getStockRegistersByStockDateRanged == nil || getLastStockRegistersByStock == nil || getStockTendencyStatByStock == nil {
		panic("nil arguments for NewStockHandler")
	}

	return &StockRegisterHandler{
		getStockRegistersByStockDateRanged,
		getLastStockRegistersByStock,
		getStockTendencyStatByStock,
	}
}
