package repository

import (
	"context"
	"testing"

	"github.com/odas0r/zet/internal/model"
	"github.com/muxit-studio/test/assert"
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
		assert.Equal(t, err, nil, "error should be nil")

		wantContent := `# Testing Zettel


`
		assert.Equal(t, zettel.Content, wantContent, "content should be empty")

		assert.Equal(t, zettel.Type, "fleet", "type should be fleet")
		assert.Equal(t, zettel.Path, "/home/odas0r/github.com/odas0r/zet/fleet/testing-zettel."+zettel.ID+".md", "path should be correct")
	})
	t.Run("can create bulk zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		repo := NewZettelRepository(db)

		zettel1 := &model.Zettel{
			Title: "Testing Zettel",
		}
		zettel2 := &model.Zettel{
			Title: "Testing Zettel 2",
		}

		err := repo.CreateBulk(context.Background(), zettel1, zettel2)
		assert.Equal(t, err, nil, "error should be nil")
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
		assert.Equal(t, err, nil, "error should be nil")

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
		assert.Equal(t, err, nil, "error should be nil")

		assert.Equal(t, len(z1.Links), 0, "z1 should not link to z2 or z3")
		assert.Equal(t, len(z2.Links), 2, "z2 should link to z3 and z1")

		err = repo.Unlink(context.Background(), z2, []*model.Zettel{z1})
		assert.Equal(t, err, nil, "error should be nil")

		assert.Equal(t, len(z2.Links), 1, "z2 should link to z1")
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
		assert.Equal(t, err, nil, "error should be nil")

		err = repo.Get(context.Background(), z1)
		assert.Equal(t, err, nil, "error should be nil")

		// check metadata
		assert.Equal(t, z1.ID, "4", "z1 should have correct id")
		assert.Equal(t, z1.Title, "Testing Zettel", "z1 should have correct title")
		assert.Equal(t, z1.Type, "fleet", "z1 should have correct type")

		// check links
		assert.Equal(t, z1.Links[0].ID, "5", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "6", "z1 should link to z3")
	})
}
