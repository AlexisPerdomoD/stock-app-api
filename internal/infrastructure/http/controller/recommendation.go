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

type RecommendationController struct {
	getRecommendationsByStockUC *usecase.GetRecommendationsByStockUseCase
}

func NewRecommendationController(getRecommendationsByStockUC *usecase.GetRecommendationsByStockUseCase) *RecommendationController {

	if getRecommendationsByStockUC == nil {
		log.Fatalln("[RecommendationController]: getRecommendationsByStockUC provided as nil")
	}

	return &RecommendationController{getRecommendationsByStockUC}
}

func (rc *RecommendationController) GetRecommendationsByStockHandler(c *gin.Context) {
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

	filters := dto.MapGetRecommendationsFilter(c)
	ctx := c.Request.Context()
	recommendations, err := rc.getRecommendationsByStockUC.Execute(ctx, filters, uint(parsedStockID))
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, recommendations)

}

func (rc *RecommendationController) SetRoutes(r *gin.Engine) {
	group := r.Group("/recommendations")
	group.Use(middleware.UserSessionMiddleware)

	group.GET("/:stockID", rc.GetRecommendationsByStockHandler)
}
