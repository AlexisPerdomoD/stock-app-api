package dto

import (
	"strconv"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

func MapGetStocksFilter(c *gin.Context) pkg.PaginationFilter {

	search := c.Query("search")
	orderBy := c.Query("orderby")

	greaterThan := c.Query("greater")
	lowerThan := c.Query("lower")
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "20")

	filters := pkg.PaginationFilter{
		SortBy: map[string]pkg.SortOrder{},

		PaginationPage: pkg.PaginationPage{
			Size: 20,
			Page: 1,
		},
		FilterBy: []pkg.FilterByItem{},
	}

	parsedGreater, err := strconv.ParseFloat(greaterThan, 64)
	if err == nil && parsedGreater > 0 {
		filters.FilterBy = append(filters.FilterBy, pkg.FilterByItem{
			Field:    "price",
			Value:    parsedGreater,
			Operator: pkg.GreaterOrEq,
		})
	}

	parsedLower, err := strconv.ParseFloat(lowerThan, 64)
	if err == nil &&
		parsedLower > 0 &&
		(parsedGreater == 0 || parsedLower > parsedGreater) {
		filters.FilterBy = append(filters.FilterBy, pkg.FilterByItem{
			Field:    "price",
			Value:    parsedLower,
			Operator: pkg.LessOrEq,
		})
	}

	if search != "" {
		filters.Search = strings.ToLower(search)
	}

	switch orderBy {
	case "tendency-asc":
		filters.SortBy["tendency"] = pkg.SortOrderAsc
		filters.SortBy["price"] = pkg.SortOrderDesc
	case "tendency-desc":
		filters.SortBy["tendency"] = pkg.SortOrderDesc
		filters.SortBy["price"] = pkg.SortOrderDesc
	case "price-asc":
		filters.SortBy["price"] = pkg.SortOrderAsc
	case "price-desc":
		filters.SortBy["price"] = pkg.SortOrderDesc
	case "ticker-asc":
		filters.SortBy["ticker"] = pkg.SortOrderAsc
	case "ticker-desc":
		filters.SortBy["ticker"] = pkg.SortOrderDesc
	case "date":
		filters.SortBy["updated_at"] = pkg.SortOrderAsc
	default:
		filters.SortBy["updated_at"] = pkg.SortOrderDesc
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
