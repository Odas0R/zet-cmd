package sqltest

import (
	"fmt"
	"testing"

	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/fs"
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

	cmd := fmt.Sprintf("goose -dir %s/migrations sqlite3 %s/zettel_test.db up", ZET_CMD_PATH, ZET_CMD_PATH)

	if err := fs.Exec(cmd); err != nil {
		t.Fatal(err)
	}

	return db
}
