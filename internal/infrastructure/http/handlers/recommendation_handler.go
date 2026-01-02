package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	getRecommendationsByStock *usecases.GetRecommendationsByStock
}

// GetRecommendationsByStockHandler godoc
// @Summary Obtener recomendaciones de un stock
// @Description Obtener recomendaciones de un stock
// @Tags recommendations
// @Produce json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar:'Bearer {token}'"
// @Param stockID path string true "ID del stock"
// @Param search query string false "Búsqueda en las recomendaciones"
// @Param groupby query string false "Agrupar por rating si se usa: 'rating'"
// @Param page query int false "Número de página. Default: 1"
// @Param size query int false "Tamaño de página. Default: 20"
// @Success 200 {object} models.RecommendationViewPaginated
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/recommendations/{stockID} [get]
func (rc *RecommendationHandler) GetRecommendationsByStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("required and not provided stockID for this service"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	parsedStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("invalid stockID provided"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	filters := mappers.MapGetRecommendationsFilter(c)
	filters.FilterBy = []pkg.FilterByItem{{
		Field:    domain.FilterByRecommendationStockID.String(),
		Operator: pkg.Equals,
		Value:    parsedStockID,
	}}

	ctx := c.Request.Context()
	recommendations, err := rc.getRecommendationsByStock.Execute(ctx, *filters)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, recommendations)

}

func NewRecommendationHandler(getRecommendationsByStockUC *usecases.GetRecommendationsByStock) *RecommendationHandler {

	if getRecommendationsByStockUC == nil {
		log.Fatalln("[RecommendationController]: getRecommendationsByStockUC provided as nil")
	}

	return &RecommendationHandler{getRecommendationsByStockUC}
}
