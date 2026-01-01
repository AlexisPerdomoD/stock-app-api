package usecases_test

import (
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/pkg"
)

func AssertApiErr(t *testing.T, err error, code int, name string) {
	t.Helper()

	apiErr, ok := err.(*pkg.ApiErr)
	if !ok {
		t.Fatalf("expected *pkg.ApiErr, got %T", err)
	}

	if apiErr.Code != code {
		t.Fatalf("expected code %d, got %d", code, apiErr.Code)
	}

	if apiErr.Name != name {
		t.Fatalf("expected name %q, got %q", name, apiErr.Name)
	}
}
