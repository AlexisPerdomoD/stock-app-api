package mappers

import (
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func MapGetRecommendationsFilter(c *gin.Context) *pkg.PaginationFilter {
	search := c.Query("search")
	groupByRating := c.Query("groupby") == "rating"

	page := c.Query("page")
	size := c.Query("size")

	filters := &pkg.PaginationFilter{
		SortBy: make([]pkg.SortByItem, 0, 2),

		PaginationPage: pkg.PaginationPage{
			Size: 20,
			Page: 1,
		},
		FilterBy: []pkg.FilterByItem{},
	}

	if search != "" {
		filters.Search = strings.ToLower(search)
	}

	filters.SortBy = append(filters.SortBy, pkg.SortByItem{
		Field: "updated_at",
		Order: pkg.SortOrderDesc,
	})

	if groupByRating {
		filters.SortBy = append(filters.SortBy, pkg.SortByItem{
			Field: "rating_to",
			Order: pkg.SortOrderDesc,
		})
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
