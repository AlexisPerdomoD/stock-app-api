package shared

import (
	"errors"
	"fmt"
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
	return &ApiErr{Detail: detail, Code: 404, Name: "Not Found"}
}

func BadRequest(detail string) error {
	return &ApiErr{Detail: detail, Code: 400, Name: "Bad Request"}
}

func Unauthorized(detail string) error {
	return &ApiErr{Detail: detail, Code: 401, Name: "Unauthorized"}
}

func Forbidden(detail string) error {
	return &ApiErr{Detail: detail, Code: 403, Name: "Forbidden"}
}

func DataBaseErr(detail string, code int) error {

	if code == 0 {
		code = 500
	}

	return &ApiErr{Detail: detail, Code: code, Name: "Database Error"}
}

func InternalServerError(detail string) error {
	return &ApiErr{Detail: detail, Code: 500, Name: "Internal Server Error"}
}

func MapHttpErr(err error) *ResponseError {
	var apiErr *ApiErr

	if err == nil || !errors.As(err, &apiErr) {
		return &ResponseError{StatusCode: 500, Name: "Unknown Error", Message: "this is a bug :("}
	}

	return &ResponseError{
		StatusCode: apiErr.Code,
		Name:       apiErr.Name,
		Message:    apiErr.Detail,
	}

}
