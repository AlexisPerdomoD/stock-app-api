package cockroachdb

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_POPULATED_STOCK_RECOMMENDATION_QUERY = `
	SELECT
		r.id 						AS id,
		r.brokerage_id				AS brokerage_id,
		r.stock_register_id			AS stock_register_id,
		r.rating_from				AS rating_from,
		r.rating_to					AS rating_to,
		r.target_from				AS target_from,
		r.target_to					AS target_to,
		r.created_at				AS created_at,

		b.name						AS brokerage_name,
		b.created_at				AS brokerage_created_at
	FROM stock_recommendations 		AS r 
	INNER JOIN brokerages			AS b
		ON b.id = r.brokerage_id
	INNER JOIN stock_registers  	AS sr
		ON sr.id = r.stock_register_id
	`

const INSERT_STOCK_RECOMMENDATION_NAMED_QUERY = `
	INSERT INTO stock_recommendations(
		brokerage_id,
		stock_register_id,
		rating_from,
		rating_to,
		target_from,
		target_to,
		created_at,
		batch_index
	) VALUES (
		:brokerage_id, 
		:stock_register_id, 
		:rating_from, 
		:rating_to, 
		:target_from, 
		:target_to, 
		:created_at,
		:batch_index
	)
	RETURNING
		id, 
		brokerage_id, 
		stock_register_id,
		rating_from,
		rating_to,
		target_from,
		target_to,
		created_at,
		batch_index
	`

type RecommendationRepository struct {
	db                          sqlx.ExtContext
	paginationFilterFieldMapper map[domain.FilterByRecommendation]FieldValidator
	paginationSortFieldMapper   map[domain.SortByRecommendation]string
}

func (r *RecommendationRepository) GetAllPaginated(
	ctx context.Context,
	filter pkg.PaginationFilter,
) (*pkg.PaginationResponse[domain.PopulatedRecommendation], error) {
	args := make([]any, 0, len(filter.FilterBy)+1)
	statement := strings.Builder{}
	statement.WriteString(GET_POPULATED_STOCK_RECOMMENDATION_QUERY)
	statement.WriteString(" WHERE 1=1")

	for _, f := range filter.FilterBy {
		field := domain.FilterByRecommendation(f.Field)
		if !field.IsValid() {
			continue
		}

		fieldValidator, ok := r.paginationFilterFieldMapper[field]
		if !ok {
			return nil, pkg.InvalidStateErr(fmt.Sprintf("invalid map of filter field %s", f.Field))
		}

		column, ok := fieldValidator.GetColumn(f.Field, f.Value)
		if !ok {
			return nil, pkg.InvalidStateErr(fmt.Sprintf("invalid map value of filter field %s %v", f.Field, f.Value))
		}

		op, err := newFilterOperator(f.Operator)
		if err != nil {
			return nil, pkg.InvalidStateErr(err.Error())
		}

		if op.RequireValue() {
			fmt.Fprintf(&statement, " AND %s %s $%d", column, op, len(args)+1)
			args = append(args, f.Value)
		} else {
			fmt.Fprintf(&statement, " AND %s %s", column, op)
		}
	}

	if filter.Search != "" {
		pos := len(args) + 1
		fmt.Fprintf(&statement, ` AND b.name ILIKE $%d ESCAPE '\'`, pos)
		args = append(args, mapILIKE(filter.Search))
	}
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count", statement.String())
	var totalSize int
	if err := r.db.QueryRowxContext(ctx, countQuery, args...).
		Scan(&totalSize); err != nil {
		return nil, err
	}

	if totalSize == 0 {
		return &pkg.PaginationResponse[domain.PopulatedRecommendation]{
			Items:      make([]domain.PopulatedRecommendation, 0),
			Page:       filter.Page,
			PageSize:   filter.Size,
			TotalSize:  0,
			TotalPages: 0,
		}, nil
	}

	statement.WriteString(" ORDER BY ")
	if len(filter.SortBy) == 0 {
		// default sort by
		statement.WriteString("r.created_at DESC")
	} else {
		orderCount := 0
		for _, item := range filter.SortBy {
			field := domain.SortByRecommendation(item.Field)
			if !field.IsValid() {
				continue
			}

			column, ok := r.paginationSortFieldMapper[field]
			if !ok {
				return nil, pkg.InvalidStateErr(fmt.Sprintf("invalid map of sort field %s", item.Field))
			}

			if orderCount > 0 {
				statement.WriteString(", ")
			}

			fmt.Fprintf(&statement, "%s %s", column, item.Order)
			orderCount++
		}
	}

	size, page := filter.GetSafeSize(), filter.GetSafePage()
	limit, offset := size, (page-1)*size
	fmt.Fprintf(&statement, " LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.QueryxContext(ctx, statement.String(), args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]domain.PopulatedRecommendation, 0, filter.Size)

	for rows.Next() {
		row := populatedStockRecommendationQueryRow{}
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}

		items = append(items, domain.PopulatedRecommendation{
			Recommendation: domain.Recommendation{
				ID:              row.ID,
				BrokerageID:     row.BrokerageID,
				StockRegisterID: row.StockRegisterID,
				RatingFrom:      row.RatingFrom,
				RatingTo:        row.RatingTo,
				TargetFrom:      row.TargetFrom,
				TargetTo:        row.TargetTo,
				CreatedAt:       row.CreatedAt,
			},
			Brokerage: domain.Brokerage{
				ID:        row.BrokerageID,
				Name:      row.BrokerageName,
				CreatedAt: row.BrokerageCreatedAt,
			},
		})
	}

	response := &pkg.PaginationResponse[domain.PopulatedRecommendation]{
		Items:      items,
		Page:       page,
		PageSize:   limit,
		TotalSize:  totalSize,
		TotalPages: int(math.Ceil(float64(totalSize) / float64(filter.Size))),
	}

	return response, rows.Err()
}

func (r *RecommendationRepository) SaveAll(
	ctx context.Context,
	recommendations []*domain.Recommendation,
) error {
	if len(recommendations) == 0 {
		return nil //no-op
	}

	args := make([]*stockRecommendationRecord, 0, len(recommendations))

	for i, recommendation := range recommendations {
		if recommendation == nil {
			return pkg.InvalidStateErr("nil pointer passed on recommendations slice")
		}

		record := &stockRecommendationRecord{
			BrokerageID:     recommendation.BrokerageID,
			StockRegisterID: recommendation.StockRegisterID,
			RatingFrom:      recommendation.RatingFrom,
			RatingTo:        recommendation.RatingTo,
			TargetFrom:      recommendation.TargetFrom,
			TargetTo:        recommendation.TargetTo,
			CreatedAt:       recommendation.CreatedAt,
			BatchIndex:      i,
		}

		args = append(args, record)
	}
	q := INSERT_STOCK_RECOMMENDATION_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	insertedCount := 0
	for rows.Next() {
		record := &stockRecommendationRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		record.MapDomain(recommendations[record.BatchIndex])
		insertedCount++
	}

	if insertedCount != len(recommendations) {
		return pkg.InvalidStateErr(fmt.Sprintf("inserted count %d != len(recommendations) %d", insertedCount, len(recommendations)))
	}

	return rows.Err()
}

func NewRecommendationRepository(db sqlx.ExtContext) *RecommendationRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}

	paginationFilterFieldMapper := map[domain.FilterByRecommendation]FieldValidator{
		domain.FilterByRecommendationStockID: {
			field:          domain.FilterByRecommendationStockID.String(),
			column:         "sr.stock_id",
			valueValidator: generateCheckType[uint64](),
		},
	}

	paginationSortFieldMapper := map[domain.SortByRecommendation]string{
		domain.SortByRecommendationCreatedAt:             "r.created_at",
		domain.SortByRecommendationStockRegisterTendency: "sr.tendency",
	}

	return &RecommendationRepository{
		db,
		paginationFilterFieldMapper,
		paginationSortFieldMapper}
}
