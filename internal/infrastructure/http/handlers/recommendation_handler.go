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

type RecommendationHandler struct {
	getRecommendationsByStockUC *usecases.GetRecommendationsByStock
}

func NewRecommendationController(getRecommendationsByStockUC *usecases.GetRecommendationsByStock) *RecommendationHandler {

	if getRecommendationsByStockUC == nil {
		log.Fatalln("[RecommendationController]: getRecommendationsByStockUC provided as nil")
	}

	return &RecommendationHandler{getRecommendationsByStockUC}
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
	recommendations, err := rc.getRecommendationsByStockUC.Execute(ctx, filters, uint(parsedStockID))
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, recommendations)

}

func (rc *RecommendationHandler) SetRoutes(r *gin.Engine) {
	group := r.Group("/recommendations")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("/:stockID", rc.GetRecommendationsByStockHandler)
}
