package handlers

import (
	"net/http"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	getMarkets *usecases.GetMarkets
}

// GetMarketsHandler godoc
// @Summary Obtener los mercados (fuentes de datos)
// @Description Obtiene todos los mercados (fuentes de datos) disponibles
// @Tags markets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Authorization Bearer token"
// @Success 200 {array}  models.MarketView
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/markets [get]
func (h *MarketHandler) GetMarketsHandler(c *gin.Context) {
	res, err := h.getMarkets.Execute(c)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func NewMarketHandler(getMarkets *usecases.GetMarkets) *MarketHandler {
	if getMarkets == nil {
		panic("nil pointer passed on getMarkets")
	}

	return &MarketHandler{getMarkets}
}
