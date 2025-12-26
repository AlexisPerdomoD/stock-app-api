package handlers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type RecommendationHandler struct {
	getRecommendationsByStock *usecases.GetRecommendationsByStock
}

func (rc *RecommendationHandler) GetRecommendationsByStockHandler(c *gin.Context) {
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

	filters := mappers.MapGetRecommendationsFilter(c)
	ctx := c.Request.Context()
	recommendations, err := rc.getRecommendationsByStock.Execute(ctx, *filters, uint(parsedStockID))
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, recommendations)

}

func (rc *RecommendationHandler) SetRoutes(r *gin.Engine) {
	group := r.Group("/recommendations")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("/:stockID", rc.GetRecommendationsByStockHandler)
}

func NewRecommendationHandler(getRecommendationsByStockUC *usecases.GetRecommendationsByStock) *RecommendationHandler {

	if getRecommendationsByStockUC == nil {
		log.Fatalln("[RecommendationController]: getRecommendationsByStockUC provided as nil")
	}

	return &RecommendationHandler{getRecommendationsByStockUC}
}
