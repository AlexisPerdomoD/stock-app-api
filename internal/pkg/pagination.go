package pkg

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

type FilterOperator string

const (
	Equals      FilterOperator = "equals"
	GreaterThan FilterOperator = "greater_than"
	GreaterOrEq FilterOperator = "greater_or_equals"
	LessThan    FilterOperator = "less_than"
	LessOrEq    FilterOperator = "less_or_equals"
	NotEquals   FilterOperator = "not_equals"
	IsNull      FilterOperator = "is_null"
	IsNotNull   FilterOperator = "is_not_null"
	Like        FilterOperator = "like"
)

type FilterByItem struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
}

type PaginationPage struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type PaginationFilter struct {
	SortBy   map[string]SortOrder `json:"sort_by"`
	FilterBy []FilterByItem       `json:"filter_by"`
	Search   string
	PaginationPage
}

type PaginationReponse[T interface{}] struct {
	Items     []T `json:"items"`
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalSize int `json:"total_size"`
}
