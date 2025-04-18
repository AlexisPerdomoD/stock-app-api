package pkg

import (
	"errors"
	"fmt"
	"net/http"
)

type ApiErr struct {
	Detail string
	Code   int
	Name   string
}

func (e *ApiErr) Error() string {

	return fmt.Sprintf("code: %d, name: %s, detail: %s", e.Code, e.Name, e.Detail)
}

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

func NotFound(detail string) error {
	return &ApiErr{Detail: detail, Code: http.StatusNotFound, Name: "Not Found"}
}

func BadRequest(detail string) error {
	return &ApiErr{Detail: detail, Code: http.StatusBadRequest, Name: "Bad Request"}
}

func Unauthorized(detail string) error {
	return &ApiErr{Detail: detail, Code: http.StatusUnauthorized, Name: "Unauthorized"}
}

func Forbidden(detail string) error {
	return &ApiErr{Detail: detail, Code: http.StatusForbidden, Name: "Forbidden"}
}

func DataBaseErr(detail string, code int) error {

	if code == 0 {
		code = http.StatusInternalServerError
	}

	return &ApiErr{Detail: detail, Code: code, Name: "Database Error"}
}

func InternalServerError(detail string) error {
	return &ApiErr{
		Detail: detail,
		Code:   http.StatusInternalServerError,
		Name:   "Internal Server Error",
	}
}

func MapHttpErr(err error) *ResponseError {
	var apiErr *ApiErr

	if err == nil {
		return &ResponseError{
			StatusCode: http.StatusInternalServerError,
			Name:       "Unknown Error",
			Message:    "Error interno desconocido",
		}
	}

	if !errors.As(err, &apiErr) {
		return &ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Name:       "Internal Server Error",
		}

	}

	return &ResponseError{
		StatusCode: apiErr.Code,
		Name:       apiErr.Name,
		Message:    apiErr.Detail,
	}

}
