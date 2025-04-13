package pkg

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type FilterOperator string

const (
	Equals      FilterOperator = "equals"
	GreaterThan FilterOperator = "greater_than"
	LessThan    FilterOperator = "less_than"
	NotEquals   FilterOperator = "not_equals"
	IsNull      FilterOperator = "is_null"
	IsNotNull   FilterOperator = "is_not_null"
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
	SortOrder SortOrder      `json:"sort_order"`
	SortBy    string         `json:"sort_by"`
	FilterBy  []FilterByItem `json:"filter_by"`
	PaginationPage
}

type PaginationReponse[T interface{}] struct {
	Items     []T `json:"items"`
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalSize int `json:"total_size"`
}
