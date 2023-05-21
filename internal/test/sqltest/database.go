package sqltest

import (
	"fmt"
	"testing"

	"github.com/odas0r/zet/pkg/database"
	"github.com/pressly/goose"
)

const (
	ZET_CMD_PATH = "/home/odas0r/github.com/odas0r/zet-cmd"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T) *database.Database {
	t.Helper()

	db := database.NewDatabase(database.NewDatabaseOptions{
		URL:                fmt.Sprintf("file:%s/zettel_test.db", ZET_CMD_PATH),
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})

	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db.DB.DB, fmt.Sprintf("%s/migrations", ZET_CMD_PATH)); err != nil {
		t.Fatal(err)
	}

	return db
}
