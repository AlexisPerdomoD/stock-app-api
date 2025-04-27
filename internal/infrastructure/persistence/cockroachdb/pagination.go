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
		case pkg.Like:
			query = query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), filter.Value)
		case pkg.Equals:
			query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		case pkg.GreaterThan:
			query = query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
		case pkg.GreaterOrEq:
			query = query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
		case pkg.LessThan:
			query = query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
		case pkg.LessOrEq:
			query = query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
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
		log.Fatalln("applyPagination: db is nil")
	}

	defaultSize := 20
	defaultPage := 1
	defaultOrder := clause.OrderByColumn{
		Column: clause.Column{Name: "created_at"},
		Desc:   true,
	}

	if filters == nil {
		filters = &pkg.PaginationFilter{
			PaginationPage: pkg.PaginationPage{
				Page: defaultPage,
				Size: defaultSize,
			},
		}
	}

	if filters.Size <= 0 {
		filters.Size = defaultSize
	}

	if filters.Page <= 0 {

		filters.Page = defaultPage
	}

	if len(filters.SortBy) == 0 {
		return query.
			Order(defaultOrder).
			Limit(filters.Size).
			Offset(filters.Size * (filters.Page - 1))
	}

	for sortBy, orderBy := range filters.SortBy {
		if !allowedSorters[sortBy] {
			continue
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, orderBy))
	}

	query = query.
		Limit(filters.Size).
		Offset(filters.Size * (filters.Page - 1))
	return query
}
