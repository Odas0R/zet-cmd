package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/muxit-studio/test/assert"
	"github.com/muxit-studio/test/require"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/repository"
	"github.com/odas0r/zet/internal/test/sqltest"
	"github.com/odas0r/zet/pkg/fs"
)

func TestNewZet(t *testing.T) {

	t.Run("create zettels -> linking -> fetching links", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(t)
		})

		zr, _ := startup(t)

		z1 := createZet(t, zr, "A title one")
		z2 := createZet(t, zr, "A title two")
		z3 := createZet(t, zr, "A title three")
		z3.WriteLine(fmt.Sprintf("This zet is linked to [[%s]]", z1.Slug))
		z3.WriteLine(fmt.Sprintf("This zet is linked to [[%s]]", z2.Slug))

		z3 = saveZet(t, zr, z3)
		require.Equal(t, len(z3.Links), 2, "z3.Links != 2")

		ids := strings.Join([]string{z3.Links[0].ID, z3.Links[1].ID}, " ")
		assert.Equal(t, strings.Contains(ids, z1.ID), true, "z1.ID not found in z3.Links")
		assert.Equal(t, strings.Contains(ids, z2.ID), true, "z2.ID not found in z3.Links")
	})

	t.Run("create zettels with fs with links -> sync -> fetching links", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(t)
		})

		zr, _ := startup(t)

		z1 := createZet(t, zr, "A title one")
		z2 := createZet(t, zr, "A title two")
		z3 := createZet(t, zr, "A title three")
		z3.WriteLine(fmt.Sprintf("This zet is linked to [[%s]]", z1.Slug))
		z3.WriteLine(fmt.Sprintf("This zet is linked to [[%s]]", z2.Slug))

		z3 = saveZet(t, zr, z3)

		err := zr.Reset(context.Background())
		require.Equal(t, err, nil, "failed to reset database")

		err = Sync(zr)
		require.Equal(t, err, nil, "failed to sync")

		err = zr.Get(context.Background(), z3)
		require.Equal(t, err, nil, "failed to fetch z3")
		require.Equal(t, len(z3.Links), 2, "z3.Links != 2")

		fmt.Println(z3.Links[0].ID, z1.ID)
		fmt.Println(z3.Links[1].ID, z2.ID)

		assert.Equal(t, z3.Links[0].ID == z1.ID, true, "z1.ID not found in z3.Links")
		assert.Equal(t, z3.Links[1].ID == z2.ID, true, "z2.ID not found in z3.Links")
	})

	t.Run("create zettels with fs with links -> save -> history", func(t *testing.T) {
		t.Cleanup(func() {
			cleanup(t)
		})

		zr, _ := startup(t)

		z1 := createZet(t, zr, "A title one")
		z2 := createZet(t, zr, "A title two")
		z3 :=	createZet(t, zr, "A title three")

		assert.Equal(t, InsertHistory(zr, z1), nil, "failed to insert history")
		assert.Equal(t, InsertHistory(zr, z2), nil, "failed to insert history")
		assert.Equal(t, InsertHistory(zr, z3), nil, "failed to insert history")

		history, err := History(zr)
		require.Equal(t, err, nil, "failed to get history")

		assert.Equal(t, len(history), 3, "history != 3")
		assert.Equal(t, history[2].Title, "A title three", "history[0].Title != 'A title three'")
		assert.Equal(t, history[1].Title, "A title two", "history[1].Title != 'A title two'")
		assert.Equal(t, history[0].Title, "A title one", "history[2].Title != 'A title one'")
	})
}

func startup(t *testing.T) (repository.ZettelRepository, *config.Config) {
	cfg := config.New("/tmp/zet-cmd")
	db := sqltest.CreateDatabase(t, cfg)
	zr := repository.NewZettelRepository(db, cfg)
	return zr, cfg
}

func createZet(t *testing.T, zr repository.ZettelRepository, title string) *model.Zettel {
	zet, err := New(zr, title)
	require.Equal(t, err, nil, "failed to create zettel")
	return zet
}

func saveZet(t *testing.T, zr repository.ZettelRepository, zet *model.Zettel) *model.Zettel {
	zettel, err := Save(zr, zet.Path)
	require.Equal(t, err, nil, "failed to save zettel")
	return zettel
}

func cleanup(t *testing.T) {
	zr, cfg := startup(t)

	err := zr.Reset(context.Background())
	require.Equal(t, err, nil, "failed to reset database")

	err = fs.RemoveAll(cfg.FleetRoot)
	require.Equal(t, err, nil, "failed to remove fleet root")

	err = fs.RemoveAll(cfg.PermanentRoot)
	require.Equal(t, err, nil, "failed to remove fleet root")
}
