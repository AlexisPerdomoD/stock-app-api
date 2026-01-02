package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/jmoiron/sqlx"
)

const GET_STOCK_QUERY = `
	SELECT 
		id, 
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin, 
		created_at, 
		updated_at 
	FROM stocks`

const GET_POPULATED_STOCK_QUERY = `
	SELECT 
		s.id									AS id,
		s.market_id								AS market_id,
		s.company_id							AS company_id,
		s.ticker								AS ticker,
		s.name									AS name,
		s.isin									AS isin,
		s.created_at							AS created_at,
		s.updated_at							AS updated_at,

		m.name									AS market_name,
		m.created_at							AS market_created_at,

		c.name									AS company_name,
		c.created_at							AS company_created_at,
		
		lsr.id									AS last_stock_register_id,
		lsr.price								AS last_stock_register_price,
		lsr.tendency							AS last_stock_tendency,
		lsr.created_at							AS last_stock_created_at

	FROM stocks 								AS s
		INNER JOIN markets 						AS m ON m.id = s.market_id
		INNER JOIN companies 					AS c ON c.id = s.company_id
		INNER JOIN LATERAL(
			SELECT 
				id,
				price, 
				tendency, 
				created_at 
			FROM stock_registers 
			WHERE stock_id = s.id 
			ORDER BY created_at DESC LIMIT 1
		) 										AS lsr ON TRUE
	`

const GET_POPULATED_STOCK_BY_USER_QUERY = `
	SELECT 
		s.id									AS id,
		s.market_id								AS market_id,
		s.company_id							AS company_id,
		s.ticker								AS ticker,
		s.name									AS name,
		s.isin									AS isin,
		s.created_at							AS created_at,
		s.updated_at							AS updated_at,

		m.name									AS market_name,
		m.created_at							AS market_created_at,

		c.name									AS company_name,
		c.created_at							AS company_created_at,
		
		lsr.id									AS last_stock_register_id,
		lsr.price								AS last_stock_register_price,
		lsr.tendency							AS last_stock_tendency,
		lsr.created_at							AS last_stock_created_at

	FROM stocks AS s
		INNER JOIN stock_users 					AS su ON su.stock_id = s.id AND su.user_id = $1
		INNER JOIN markets 						AS m ON m.id = s.market_id
		INNER JOIN companies 					AS c ON c.id = s.company_id
		INNER JOIN LATERAL(
			SELECT 
				id,
				price, 
				tendency, 
				created_at 
			FROM stock_registers 
			WHERE stock_id = s.id 
			ORDER BY created_at DESC LIMIT 1
		) 										AS lsr ON TRUE
	`

const INSERT_STOCK_QUERY = `
	INSERT INTO stocks(
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		market_id,
		company_id,
		ticker,
		name,
		isin,
		created_at,
		updated_at`

const INSERT_STOCK_NAMED_QUERY = `
	INSERT INTO stocks(
		market_id, 
		company_id, 
		ticker, 
		name, 
		isin,
		batch_index
	) VALUES (
		:market_id, 
		:company_id, 
		:ticker, 
		:name, 
		:isin, 
		:batch_index
	)
	RETURNING
		id,
		market_id,
		company_id,
		ticker,
		name,
		isin,
		created_at,
		updated_at,
		batch_index`

type StockRepository struct {
	db sqlx.ExtContext

	filterByFieldMap map[domain.FilterByStock]FieldValidator
	orderByFieldMap  map[domain.SortByStock]string
}

func (r *StockRepository) GetByID(ctx context.Context, stockID uint64) (*domain.Stock, error) {
	record := &stockRecord{}
	q := GET_STOCK_QUERY + " WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, q, stockID).
		StructScan(record); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return record.ToDomain(), nil
}

func (r *StockRepository) paginate(
	ctx context.Context,
	filter pkg.PaginationFilter,
	userID *uint64,
) (*pkg.PaginationResponse[domain.PopulatedStock], error) {

	var q string
	args := make([]any, 0, len(filter.FilterBy)+1)
	statement := strings.Builder{}
	statement.WriteString(" WHERE 1=1")

	if userID != nil {
		q = GET_POPULATED_STOCK_BY_USER_QUERY
		args = append(args, userID)
	} else {
		q = GET_POPULATED_STOCK_QUERY
	}

	for _, f := range filter.FilterBy {
		op, err := newFilterOperator(f.Operator)
		if err != nil {
			return nil, pkg.InvalidStateErr(err.Error())
		}

		field := domain.FilterByStock(f.Field)
		if !field.IsValid() {
			continue
		}

		fieldValidator, ok := r.filterByFieldMap[field]
		if !ok {
			continue
		}

		column, ok := fieldValidator.GetColumn(f.Field, f.Value)
		if !ok {
			continue
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
		fmt.Fprintf(&statement,
			` AND (
				s.ticker ILIKE $%d ESCAPE '\'
				OR s.name ILIKE $%d ESCAPE '\'
				OR c.name ILIKE $%d ESCAPE '\'
			)`,
			pos, pos, pos)
		args = append(args, mapILIKE(filter.Search))
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count", q+statement.String())

	statement.WriteString(" ORDER BY ")
	if len(filter.SortBy) == 0 {
		// default sort by
		statement.WriteString("lsr.created_at DESC")
	} else {
		orderCount := 0
		for _, item := range filter.SortBy {
			field := domain.SortByStock(item.Field)
			if !field.IsValid() {
				continue
			}

			column, ok := r.orderByFieldMap[field]
			if !ok {
				return nil, pkg.InvalidStateErr(fmt.Sprintf("invalid sort field %s", item.Field))
			}

			if orderCount > 0 {
				statement.WriteString(", ")
			}

			fmt.Fprintf(&statement, "%s %s", column, item.Order)
		}
	}

	var totalRecords int
	if err := r.db.QueryRowxContext(ctx, countQuery, args...).
		Scan(&totalRecords); err != nil {
		return nil, err
	}
	size, page := filter.GetSafeSize(), filter.GetSafePage()
	totalPages := int(math.Ceil(float64(totalRecords) / float64(size)))

	if totalRecords == 0 {
		return &pkg.PaginationResponse[domain.PopulatedStock]{
			Items:      make([]domain.PopulatedStock, 0),
			Page:       page,
			PageSize:   size,
			TotalSize:  0,
			TotalPages: 0,
		}, nil
	}

	limit, offset := size, (page-1)*size
	fmt.Fprintf(&statement, " LIMIT %d OFFSET %d", limit, offset)
	rows, err := r.db.QueryContext(ctx, q+statement.String(), args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]domain.PopulatedStock, 0, filter.Size)
	for rows.Next() {
		stockRecord := &stockRecord{}
		marketRecord := &marketRecord{}
		companyRecord := &companyRecord{}
		stockRegisterRecord := &stockRegisterRecord{}

		if err := rows.Scan(
			&stockRecord.ID,
			&marketRecord.ID,
			&companyRecord.ID,
			&stockRecord.Ticker,
			&stockRecord.Name,
			&stockRecord.Isin,
			&stockRecord.CreatedAt,
			&stockRecord.UpdatedAt,

			&marketRecord.Name,
			&marketRecord.CreatedAt,

			&companyRecord.Name,
			&companyRecord.CreatedAt,

			&stockRegisterRecord.ID,
			&stockRegisterRecord.Price,
			&stockRegisterRecord.Tendency,
			&stockRegisterRecord.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, domain.PopulatedStock{
			Stock: *stockRecord.ToDomain(),
			Market: domain.Market{
				ID:        stockRecord.MarketID,
				Name:      marketRecord.Name,
				CreatedAt: marketRecord.CreatedAt,
			},
			Company: domain.Company{
				ID:        stockRecord.CompanyID,
				MarketID:  stockRecord.MarketID,
				Name:      companyRecord.Name,
				CreatedAt: companyRecord.CreatedAt,
			},
			LastRegister: domain.StockRegister{
				ID:        stockRegisterRecord.ID,
				StockID:   stockRecord.ID,
				Price:     stockRegisterRecord.Price,
				Tendency:  stockRegisterRecord.Tendency,
				CreatedAt: stockRegisterRecord.CreatedAt,
			},
		})
	}

	response := &pkg.PaginationResponse[domain.PopulatedStock]{
		Items:      items,
		Page:       page,
		PageSize:   size,
		TotalSize:  totalRecords,
		TotalPages: totalPages,
	}
	return response, nil
}

func (r *StockRepository) GetAllPaginated(ctx context.Context, filter pkg.PaginationFilter) (*pkg.PaginationResponse[domain.PopulatedStock], error) {
	return r.paginate(ctx, filter, nil)
}

func (r *StockRepository) GetAllPaginatedByUser(ctx context.Context, filter pkg.PaginationFilter, userID uint64) (*pkg.PaginationResponse[domain.PopulatedStock], error) {
	return r.paginate(ctx, filter, &userID)
}

func (r *StockRepository) GetByStockCompanySearchParams(
	ctx context.Context,
	searchParams []domain.StockCompanySearchParam,
) (map[domain.StockCompanySearchParam]*domain.Stock, error) {
	results := make(map[domain.StockCompanySearchParam]*domain.Stock)
	for _, key := range searchParams {
		results[key] = nil
	}

	if len(searchParams) == 0 {
		return results, nil
	}

	args := make([]any, 0, len(results)*3)
	where := strings.Builder{}
	where.WriteString(" WHERE (market_id, company_id, ticker) IN  (")

	i := 0
	first := true
	for key := range results {
		if !first {
			where.WriteString(", ")
		}

		fmt.Fprintf(&where, "($%d, $%d, $%d)", i+1, i+2, i+3)
		args = append(args, key.MarketID, key.CompanyID, key.StockTicker)

		i += 3
		first = false
	}

	where.WriteString(")")
	q := GET_STOCK_QUERY + where.String()
	rows, err := r.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRecord{}
		if err := rows.StructScan(record); err != nil {
			return nil, err
		}

		key := domain.StockCompanySearchParam{
			MarketID:    record.MarketID,
			CompanyID:   record.CompanyID,
			StockTicker: record.Ticker,
		}

		results[key] = record.ToDomain()
	}

	return results, rows.Err()
}

func (r *StockRepository) Save(ctx context.Context, stock *domain.Stock) error {
	if stock == nil {
		return pkg.InvalidStateErr("nil pointer passed on stock")
	}

	var isin, name sql.NullString
	if stock.Name != nil {
		name.String = *stock.Name
		name.Valid = true
	}

	if stock.Isin != nil {
		isin.String = *stock.Isin
		isin.Valid = true
	}
	record := &stockRecord{}
	if err := r.db.QueryRowxContext(
		ctx,
		INSERT_STOCK_QUERY,
		stock.MarketID,
		stock.CompanyID,
		stock.Ticker,
		name,
		isin,
	).StructScan(record); err != nil {
		return err
	}

	record.MapDomain(stock)
	return nil
}

func (r *StockRepository) SaveAll(ctx context.Context, stocks []*domain.Stock) error {
	if len(stocks) == 0 {
		return nil //no-op
	}

	args := make([]stockRecord, 0, len(stocks))
	for i, stock := range stocks {
		if stock == nil {
			return pkg.InvalidStateErr("nil pointer passed on stock slice")
		}

		arg := stockRecord{
			MarketID:   stock.MarketID,
			CompanyID:  stock.CompanyID,
			Ticker:     stock.Ticker,
			Isin:       sql.NullString{},
			Name:       sql.NullString{},
			BatchIndex: i,
		}

		if stock.Name != nil {
			arg.Name.String = *stock.Name
			arg.Name.Valid = true
		}

		if stock.Isin != nil {
			arg.Isin.String = *stock.Isin
			arg.Isin.Valid = true
		}

		args = append(args, arg)
	}
	q := INSERT_STOCK_NAMED_QUERY
	rows, err := sqlx.NamedQueryContext(ctx, r.db, q, args)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		record := &stockRecord{}
		if err := rows.StructScan(record); err != nil {
			return err
		}

		record.MapDomain(stocks[record.BatchIndex])
	}

	return rows.Err()
}

func (r *StockRepository) Update(ctx context.Context, stock domain.StockUpdates) error {
	return pkg.InternalServerError("not implemented")
}

func NewStockRepository(db sqlx.ExtContext) *StockRepository {
	if db == nil {
		panic("db sqlx.ExtContext is nil")
	}
	filterByFieldMap := make(map[domain.FilterByStock]FieldValidator)
	filterByFieldMap[domain.FilterByStockPrice] = FieldValidator{
		field:          domain.FilterByStockPrice.String(),
		column:         "lsr.price",
		valueValidator: generateCheckType[float64](),
	}

	orderByFieldMap := make(map[domain.SortByStock]string)
	orderByFieldMap[domain.SortByStockTendency] = "lsr.tendency"
	orderByFieldMap[domain.SortByStockPrice] = "lsr.price"
	orderByFieldMap[domain.SortByStockTicker] = "s.ticker"
	orderByFieldMap[domain.SortByStockDate] = "lsr.created_at"
	return &StockRepository{db, filterByFieldMap, orderByFieldMap}
}
