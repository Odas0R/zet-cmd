package sqltest

import (
	"fmt"
	"testing"

	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/pkg/database"
	"github.com/pressly/goose"
)

const migrationsPath = "/home/odas0r/github.com/odas0r/zet-cmd/migrations"

// CreateDatabase for testing.
func CreateDatabase(t *testing.T, cfg *config.Config) *database.Database {
	t.Helper()

	db := database.NewDatabase(database.NewDatabaseOptions{
		URL:                fmt.Sprintf("file:%s/zettel_test.db", cfg.Root),
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})

	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db.DB.DB, migrationsPath); err != nil {
		t.Fatal(err)
	}

	return db
}
