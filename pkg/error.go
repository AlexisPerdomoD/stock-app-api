package pkg

import (
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

func InvalidStateErr(detail string) *ApiErr {
	return &ApiErr{
		Detail: detail,
		Code:   http.StatusInternalServerError,
		Name:   "Invalid State Error",
	}
}

func InternalServerError(detail string) error {
	return &ApiErr{
		Detail: detail,
		Code:   http.StatusInternalServerError,
		Name:   "Internal Server Error",
	}
}
