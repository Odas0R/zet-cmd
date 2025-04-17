package sqltest

import (
	"fmt"
	"os"
	"testing"

	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/pkg/database"
	"github.com/pressly/goose"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T, cfg *config.Config) *database.Database {
	home := os.Getenv("HOME")

	migrationsPath := fmt.Sprintf("%s/github.com/odas0r/zet-cmd/migrations", home)

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
