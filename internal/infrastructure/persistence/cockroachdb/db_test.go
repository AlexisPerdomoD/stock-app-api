package cockroachdb_test

import (
	"os"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
)

func TestMain(m *testing.M) {

	mustSet := func(key, value string) {
		if err := os.Setenv(key, value); err != nil {
			panic("[SetTestSetup]:failed to set env var " + key + ": " + err.Error())
		}
	}

	mustSet("CR_HOST", "localhost")
	mustSet("CR_PORT", "26258")
	mustSet("CR_USER", "root")
	mustSet("CR_PASSWORD", "")
	mustSet("CR_DB", "defaultdb")
	mustSet("CR_SSL", "disable")
	mustSet("CR_RUN_MIGRATE", "TRUE")

	db := cockroachdb.NewDB()
	if err := cockroachdb.Migrate(db); err != nil {
		panic("[SetTestSetup]: failed when migrating on db: " + err.Error())
	}

	code := m.Run()
	os.Exit(code)
}
