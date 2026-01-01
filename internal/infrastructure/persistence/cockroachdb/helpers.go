package cockroachdb

import (
	"errors"
	"strings"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

type FieldValidator struct {
	field          string
	column         string
	valueValidator func(val any) bool
}

func (f *FieldValidator) GetColumn(field string, value any) (string, bool) {
	if f.field != field {
		return "", false
	}

	if !f.valueValidator(value) {
		return "", false
	}

	return f.column, true
}

type FilterOperator string

const (
	Equals      FilterOperator = "="
	GreaterThan FilterOperator = ">"
	GreaterOrEq FilterOperator = ">="
	LessThan    FilterOperator = "<"
	LessOrEq    FilterOperator = "<="
	IsNull      FilterOperator = "IS NULL"
	IsNotNull   FilterOperator = "IS NOT NULL"
	Like        FilterOperator = "LIKE"
)

func (f FilterOperator) String() string {
	return string(f)
}

func (f FilterOperator) RequireValue() bool {
	return f != IsNull && f != IsNotNull
}

func newFilterOperator(operator pkg.FilterOperator) (FilterOperator, error) {
	var op FilterOperator
	switch operator {
	case pkg.Equals:
		op = Equals
	case pkg.GreaterThan:
		op = GreaterThan
	case pkg.GreaterOrEq:
		op = GreaterOrEq
	case pkg.LessThan:
		op = LessThan
	case pkg.LessOrEq:
		op = LessOrEq
	case pkg.IsNull:
		op = IsNull
	case pkg.IsNotNull:
		op = IsNotNull
	case pkg.Like:
		op = Like
	default:
		return "", errors.New("ilegal filter operator found")
	}

	return op, nil

}

func generateCheckType[T any]() func(val any) bool {
	return func(val any) bool {
		_, ok := val.(T)
		return ok
	}
}

func mapILIKE(search string) string {
	search = strings.ReplaceAll(search, `\`, `\\`)
	search = strings.ReplaceAll(search, `%`, `\%`)
	search = strings.ReplaceAll(search, `_`, `\_`)
	return "%" + search + "%"
}
