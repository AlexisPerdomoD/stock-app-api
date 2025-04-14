package cockroachdb

import (
	"fmt"
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applyFilters(
	query *gorm.DB,
	filters []pkg.FilterByItem,
	allowedFilters map[string]bool,
) *gorm.DB {
	if query == nil {
		log.Fatalln("bad impl: db not provided when calling ApplyFilters helper")
	}

	for _, filter := range filters {

		if !allowedFilters[filter.Field] {
			continue
		}

		switch filter.Operator {
		case pkg.Equals:
			query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		case pkg.GreaterThan:
			query = query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
		case pkg.LessThan:
			query = query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
		case pkg.NotEquals:
			query = query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
		case pkg.IsNull:
			query = query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
		case pkg.IsNotNull:
			query = query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
		default:
			continue
		}
	}

	return query
}

func applyPagination(query *gorm.DB,
	filters *pkg.PaginationFilter,
	allowedSorters map[string]bool,
) *gorm.DB {

	if query == nil {
		log.Fatalln("bad impl: db not provided when calling ApplyFilters helper")
	}

	size := 20
	page := 1
	sortBy := "created_at"
	sortOrder := pkg.Desc

	if filters == nil {
		return query.
			Limit(size).
			Offset(size * (page - 1)).
			Order(clause.OrderByColumn{
				Column: clause.Column{Name: sortBy},
				Desc:   true,
			})
	}

	if allowedSorters[filters.SortBy] {
		sortBy = filters.SortBy
	}

	if filters.SortOrder == pkg.Asc {
		sortOrder = pkg.Asc
	}

	if filters.Size > 0 {
		size = filters.Size
	}

	if filters.Page > 1 {
		page = filters.Page
	}

	return query.Limit(size).
		Offset(size * (page - 1)).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortBy},
			Desc:   sortOrder == pkg.Desc,
		})
}
