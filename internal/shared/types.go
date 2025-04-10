package shared

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type PaginationReponse[T interface{}] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginationFilter struct {
	SortOrder SortOrder `json:"sort_order"`
	SortBy    string    `json:"sort_by"`

	FilterBy []struct {
		Field string      `json:"field"`
		Value interface{} `json:"value"`
	} `json:"filter_by"`

	Page int `json:"page"`
	Size int `json:"size"`
}
