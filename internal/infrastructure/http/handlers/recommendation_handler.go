package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	getRecommendationsByStock *usecases.GetRecommendationsByStock
}

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
	ctx := c.Request.Context()
	recommendations, err := rc.getRecommendationsByStock.Execute(ctx, *filters, parsedStockID)
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
