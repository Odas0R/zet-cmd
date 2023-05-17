package repository

import (
	"context"
	"testing"

	"github.com/muxit-studio/test/assert"
	"github.com/muxit-studio/test/require"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/test/sqltest"
)

func TestZettelRepository_Create(t *testing.T) {
	t.Run("can create a zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		zettel := &model.Zettel{
			Title: "Testing Zettel",
		}

		err := repo.Create(context.Background(), zettel)
		require.Equal(t, err, nil, "failed to create zettel")

		assert.Equal(t, zettel.Lines[0], "# Testing Zettel", "first line should be the title")
		assert.Equal(t, zettel.Lines[1], "", "second line should be a empty line")
		assert.Equal(t, zettel.Lines[2], "", "third line should be a empty line")

		assert.Equal(t, zettel.Type, "fleet", "type should be fleet")
		assert.Equal(t, zettel.Path, "/home/odas0r/github.com/odas0r/zet/fleet/testing-zettel."+zettel.ID+".md", "path should be correct")
	})
	t.Run("can create bulk zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		z1 := &model.Zettel{
			Title: "Testing Zettel",
		}
		z2 := &model.Zettel{
			Title: "Testing Zettel 2",
		}
		z3 := &model.Zettel{
			Title: "Testing Zettel 3",
		}

		var zettels []*model.Zettel
		zettels = append(zettels, z1, z2, z3)

		err := repo.CreateBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettels in bulk mode")

		err = repo.Get(context.Background(), z1)
		require.Equal(t, err, nil, "failed to get zettel 1")

		assert.Equal(t, z1.Lines[0], "# Testing Zettel", "first line should be the title")
		assert.Equal(t, z1.Lines[1], "", "second line should be a empty line")
		assert.Equal(t, z1.Lines[2], "", "third line should be a empty line")

		assert.Equal(t, z1.Type, "fleet", "type should be fleet")
		assert.Equal(t, z1.Path, "/home/odas0r/github.com/odas0r/zet/fleet/testing-zettel."+z1.ID+".md", "path should be correct")

	})
}

func TestZettelRepository_Link(t *testing.T) {
	t.Run("can link different zettels", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		z1 := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
		}
		z2 := &model.Zettel{
			ID:    "2",
			Title: "Testing Zettel 2",
		}
		z3 := &model.Zettel{
			ID:    "3",
			Title: "Testing Zettel 3",
		}

		repo.Create(context.Background(), z1)
		repo.Create(context.Background(), z2)
		repo.Create(context.Background(), z3)

		err := repo.Link(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to link zettels")

		assert.Equal(t, z1.Links[0].ID, "2", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "3", "z1 should link to z3")
	})

	t.Run("can unlink different zettels", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		z1 := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
		}
		z2 := &model.Zettel{
			ID:    "2",
			Title: "Testing Zettel 2",
		}
		z3 := &model.Zettel{
			ID:    "3",
			Title: "Testing Zettel 3",
		}

		repo.Create(context.Background(), z1)
		repo.Create(context.Background(), z2)
		repo.Create(context.Background(), z3)

		repo.Link(context.Background(), z1, []*model.Zettel{z2, z3})
		repo.Link(context.Background(), z2, []*model.Zettel{z1, z3})

		err := repo.Unlink(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to unlink zettels")

		assert.Equal(t, len(z1.Links), 0, "z1 should not link to z2 or z3")
		assert.Equal(t, len(z2.Links), 2, "z2 should link to z3 and z1")

		err = repo.Unlink(context.Background(), z2, []*model.Zettel{z1})
		assert.Equal(t, err, nil, "error should be nil")

		assert.Equal(t, len(z2.Links), 1, "z2 should link to z1")
	})

	t.Run("can link bulk", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		z1 := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
		}
		z2 := &model.Zettel{
			ID:    "2",
			Title: "Testing Zettel 2",
		}
		z3 := &model.Zettel{
			ID:    "3",
			Title: "Testing Zettel 3",
		}

		repo.Create(context.Background(), z1)
		repo.Create(context.Background(), z2)
		repo.Create(context.Background(), z3)

		l1 := &model.Link{
			From: z1.ID,
			To:   z2.ID,
		}
		l2 := &model.Link{
			From: z1.ID,
			To:   z3.ID,
		}
		l3 := &model.Link{
			From: z2.ID,
			To:   z1.ID,
		}
		l4 := &model.Link{
			From: z2.ID,
			To:   z3.ID,
		}

		links := []*model.Link{l1, l2, l3, l4}

		err := repo.LinkBulk(context.Background(), links...)
		require.Equal(t, err, nil, "failed to bulk link zettels")

		repo.Get(context.Background(), z1)
		assert.Equal(t, z1.Links[0].ID, "2", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "3", "z1 should link to z3")
	})
}

func TestZettelRepository_Get(t *testing.T) {
	t.Run("can get a zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		z1 := &model.Zettel{
			ID:    "4",
			Title: "Testing Zettel",
		}
		z2 := &model.Zettel{
			ID:    "5",
			Title: "Testing Zettel 2",
		}
		z3 := &model.Zettel{
			ID:    "6",
			Title: "Testing Zettel 3",
		}

		repo.Create(context.Background(), z1)
		repo.Create(context.Background(), z2)
		repo.Create(context.Background(), z3)

		err := repo.Link(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to link zettels")

		err = repo.Get(context.Background(), z1)
		require.Equal(t, err, nil, "failed to get zettel")

		// check metadata
		assert.Equal(t, z1.ID, "4", "z1 should have correct id")
		assert.Equal(t, z1.Title, "Testing Zettel", "z1 should have correct title")
		assert.Equal(t, z1.Type, "fleet", "z1 should have correct type")
		assert.Equal(t, z1.Links[0].ID, "5", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "6", "z1 should link to z3")
	})
}
