package pkg

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type FilterByItem struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type PaginationFilter struct {
	SortOrder SortOrder      `json:"sort_order"`
	SortBy    string         `json:"sort_by"`
	FilterBy  []FilterByItem `json:"filter_by"`
	Page      int            `json:"page"`
	Size      int            `json:"size"`
}

type PaginationReponse[T interface{}] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}
