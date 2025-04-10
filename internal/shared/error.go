package shared

type ErrType string

const (
	ExternalErr     ErrType = "external_error"
	ErrTypeInternal ErrType = "internal_server_error"
	BadRequest      ErrType = "bad_request"
	Unauthorized    ErrType = "unauthorized"
	Forbidden       ErrType = "forbidden"
	NotFound        ErrType = "not_found"
)

type ApiErr struct {
	Code    int
	Message string
	Details string
	Type    ErrType
}

type ReponseError struct {
	StatusCode int         `json:"status_code"`
	Name       string      `json:"name"`
	Message    string      `json:"message"`
	Error      interface{} `json:"data"`
}
