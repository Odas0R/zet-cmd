package sqltest

import (
	"fmt"
	"testing"

	"github.com/odas0r/zet/pkg/database"
	"github.com/pressly/goose"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T, root string) *database.Database {
	t.Helper()

	db := database.NewDatabase(database.NewDatabaseOptions{
		URL:                fmt.Sprintf("file:%s/zettel_test.db", root),
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})

	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db.DB.DB, fmt.Sprintf("%s/migrations", root)); err != nil {
		t.Fatal(err)
	}

	return db
}
