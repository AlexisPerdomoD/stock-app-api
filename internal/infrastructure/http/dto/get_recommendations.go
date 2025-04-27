package dto

import (
	"strconv"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

func MapGetRecommendationsFilter(c *gin.Context) pkg.PaginationFilter {
	search := c.Query("search")
	groupByRating := c.Query("groupby") == "rating"

	page := c.Query("page")
	size := c.Query("size")

	filters := pkg.PaginationFilter{
		SortBy: map[string]pkg.SortOrder{},

		PaginationPage: pkg.PaginationPage{
			Size: 20,
			Page: 1,
		},
		FilterBy: []pkg.FilterByItem{},
	}

	if search != "" {
		filters.Search = strings.ToLower(search)
	}

	filters.SortBy["updated_at"] = pkg.SortOrderDesc

	if groupByRating {
		filters.SortBy["rating_to"] = pkg.SortOrderDesc
	}

	parsedSize, err := strconv.Atoi(size)
	if err == nil && parsedSize > 0 && parsedSize < 100 {
		filters.Size = parsedSize
	}

	parsedPage, err := strconv.Atoi(page)
	if err == nil {
		filters.Page = parsedPage
	}

	return filters
}
