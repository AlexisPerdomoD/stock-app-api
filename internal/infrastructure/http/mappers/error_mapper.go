package mappers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

func MapHttpErr(err error) *ResponseError {
	var apiErr *pkg.ApiErr

	if err == nil {
		return &ResponseError{
			StatusCode: http.StatusInternalServerError,
			Name:       "Unknown Error",
			Message:    "Error interno desconocido",
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &ResponseError{
			StatusCode: http.StatusNotFound,
			Name:       "Not Found",
			Message:    "Recurso no encontrado",
		}
	}

	if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return &ResponseError{
			StatusCode: http.StatusConflict,
			Name:       "Conflict",
			Message:    "Registro duplicado",
		}
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return &ResponseError{
			StatusCode: http.StatusConflict,
			Name:       "DataBase Err",
			Message:    "Accion no valida",
		}
	}

	if strings.Contains(err.Error(), "violates check constraint") {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Name:       "Invalid Data",
			Message:    "Datos inválidos según reglas de la base de datos",
		}
	}

	if strings.Contains(err.Error(), "null value in column") || strings.Contains(err.Error(), "violates not-null constraint") {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Name:       "Missing Field",
			Message:    "Falta un campo obligatorio",
		}
	}

	if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "could not connect to server") {
		return &ResponseError{
			StatusCode: http.StatusServiceUnavailable,
			Name:       "DB Connection Error",
			Message:    "No se pudo conectar con la base de datos",
		}
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) || strings.Contains(err.Error(), "violates foreign key constraint") {
		return &ResponseError{
			StatusCode: http.StatusConflict,
			Name:       "Database Error",
			Message:    "Violación de integridad referencial",
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

func MapValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				errors[field] = "Este campo es obligatorio"
			case "email":
				errors[field] = "Debe ser un email válido"
			case "min":
				errors[field] = "Debe tener al menos " + e.Param() + " caracteres"
			case "max":
				errors[field] = "No puede tener más de " + e.Param() + " caracteres"
			default:
				errors[field] = "Campo inválido: " + e.Tag()
			}
		}
	} else {
		errors["general"] = "Invalid Format"
	}

	return errors
}
