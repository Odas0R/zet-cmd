package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/muxit-studio/test/assert"
	"github.com/muxit-studio/test/require"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/test/sqltest"
)

var cfg = config.New("/tmp/zet-cmd")

func TestZettelRepository_Get(t *testing.T) {
	t.Run("can get a zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		repo.Save(context.Background(), z1)
		repo.Save(context.Background(), z2)
		repo.Save(context.Background(), z3)

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

func TestZettelRepository_Create(t *testing.T) {
	t.Run("can create a zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		zettel := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
		}

		createZettel(t, repo, zettel)

		err := repo.Get(context.Background(), zettel)
		require.Equal(t, err, nil, "failed to get the zettel from db")

		assert.Equal(t, zettel.Lines[0], "# Testing Zettel", "first line should be the title")
		assert.Equal(t, zettel.Lines[1], "", "second line should be a empty line")
		assert.Equal(t, zettel.Lines[2], "", "third line should be a empty line")
		assert.Equal(t, zettel.Type, "fleet", "type should be fleet")
		assert.Equal(t, zettel.Path, "/tmp/zet-cmd/fleet/testing-zettel."+zettel.ID+".md", "path should be correct")
	})
	t.Run("can create bulk zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		z1 := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
			Content: `# Testing Zettel aiosjfoiasj foiaasfiajsof  oasjdf oi
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
tempor invidunt ut labore et dolore magna aliquyam
`,
			Path: "/tmp/zet-cmd/fleet/testing-zettel.1.md",
			Type: "fleet",
		}
		z2 := &model.Zettel{
			ID:    "2",
			Title: "Testing Zettel",
			Content: `# Testing Zettel aiosjfoiasj foia
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
tempor invidunt ut labore et dolore magna aliquyam
`,
			Path: "/tmp/zet-cmd/fleet/testing-zettel.2.md",
			Type: "fleet",
		}
		z3 := &model.Zettel{
			ID:    "3",
			Title: "Testing Zettel",
			Content: `# Testing Zettel aiosjfo
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
tempor invidunt ut labore et dolore magna aliquyam
`,
			Path: "/tmp/zet-cmd/fleet/testing-zettel.3.md",
			Type: "fleet",
		}

		var zettels []*model.Zettel
		zettels = append(zettels, z1, z2, z3)

		err := repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettels in bulk mode")

		//
		// Fetching zettels
		//

		updateZettel(t, repo, z1)

		assert.Equal(t, z1.Title, "Testing Zettel", "z1 title should be correct")
		assert.Equal(t, z1.Type, "fleet", "z1 type should be fleet")
		assert.Equal(t, z1.Path, "/tmp/zet-cmd/fleet/testing-zettel."+z1.ID+".md", "z1 path should be correct")
		assert.Equal(t, z1.Lines[0], "# Testing Zettel aiosjfoiasj foiaasfiajsof  oasjdf oi", "z1 first line should be correct")

		updateZettel(t, repo, z2)

		assert.Equal(t, z2.Title, "Testing Zettel", "z2 title should be correct")
		assert.Equal(t, z2.Type, "fleet", "z2 type should be fleet")
		assert.Equal(t, z2.Path, "/tmp/zet-cmd/fleet/testing-zettel."+z2.ID+".md", "z2 path should be correct")
		assert.Equal(t, z2.Lines[0], "# Testing Zettel aiosjfoiasj foia", "z2 line 1 should be correct")

		updateZettel(t, repo, z3)

		assert.Equal(t, z3.Title, "Testing Zettel", "z3 title should be correct")
		assert.Equal(t, z3.Type, "fleet", "z3 type should be fleet")
		assert.Equal(t, z3.Path, "/tmp/zet-cmd/fleet/testing-zettel."+z3.ID+".md", "z3 path should be correct")
		assert.Equal(t, z3.Lines[0], "# Testing Zettel aiosjfo", "z3 line 1 should be correct")
	})
}

func TestZettelRepository_Link(t *testing.T) {
	t.Run("can link different zettels", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		createZettel(t, repo, z1)
		createZettel(t, repo, z2)
		createZettel(t, repo, z3)

		err := repo.Link(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to link zettels")

		assert.Equal(t, len(z1.Links), 2, "z1 should have 2 links")
		assert.Equal(t, z1.Links[0].ID, "2", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "3", "z1 should link to z3")
	})

	t.Run("can unlink different zettels", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		createZettel(t, repo, z1)
		createZettel(t, repo, z2)
		createZettel(t, repo, z3)

		err := repo.Link(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to link zettels")

		err = repo.Link(context.Background(), z2, []*model.Zettel{z1, z3})
		require.Equal(t, err, nil, "failed to link zettels")

		err = repo.Unlink(context.Background(), z1, []*model.Zettel{z2, z3})
		require.Equal(t, err, nil, "failed to unlink zettels")

		assert.Equal(t, len(z1.Links), 0, "z1 should not link to z2 or z3")
		assert.Equal(t, len(z2.Links), 2, "z2 should link to z3 and z1")

		err = repo.Unlink(context.Background(), z2, []*model.Zettel{z1})
		assert.Equal(t, err, nil, "error should be nil")
		assert.Equal(t, len(z2.Links), 1, "z2 should link to z1")
	})

	t.Run("can get backlinks", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		createZettel(t, repo, z1)
		createZettel(t, repo, z2)
		createZettel(t, repo, z3)

		links := []*model.Link{
			{
				From: z1.ID,
				To:   z2.ID,
			},
			{
				From: z1.ID,
				To:   z3.ID,
			},
			{
				From: z2.ID,
				To:   z1.ID,
			},
			{
				From: z2.ID,
				To:   z3.ID,
			},
		}

		err := repo.LinkBulk(context.Background(), links...)
		require.Equal(t, err, nil, "failed to bulk link zettels")

		backlinks, err := repo.Backlinks(context.Background(), z3)
		require.Equal(t, err, nil, "failed to query the backlinks")

		ids := strings.Join(([]string{backlinks[0].ID, backlinks[1].ID}), " ")

		assert.Equal(t, len(backlinks), 2, "z3 should have 2 backlinks")
		assert.Equal(t, strings.Contains(ids, "1"), true, "z3 should have a backlink to z1")
		assert.Equal(t, strings.Contains(ids, "2"), true, "z3 should have a backlink to z2")
	})

	t.Run("can link bulk", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		createZettel(t, repo, z1)
		createZettel(t, repo, z2)
		createZettel(t, repo, z3)

		links := []*model.Link{
			{
				From: z1.ID,
				To:   z2.ID,
			},
			{
				From: z1.ID,
				To:   z3.ID,
			},
			{
				From: z2.ID,
				To:   z1.ID,
			},
			{
				From: z2.ID,
				To:   z3.ID,
			},
		}

		err := repo.LinkBulk(context.Background(), links...)
		require.Equal(t, err, nil, "failed to bulk link zettels")

		repo.Get(context.Background(), z1)
		assert.Equal(t, z1.Links[0].ID, "2", "z1 should link to z2")
		assert.Equal(t, z1.Links[1].ID, "3", "z1 should link to z3")
	})
}

func TestZettelRepository_Remove(t *testing.T) {
	t.Run("can remove a zettel", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

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

		createZettel(t, repo, z1)
		createZettel(t, repo, z2)
		createZettel(t, repo, z3)

		links := []*model.Link{
			{
				From: z3.ID,
				To:   z1.ID,
			},
			{
				From: z2.ID,
				To:   z1.ID,
			},
		}

		err := repo.LinkBulk(context.Background(), links...)
		require.Equal(t, err, nil, "failed to bulk link zettels")

		err = repo.Remove(context.Background(), z1)
		require.Equal(t, err, nil, "failed to remove zettel")

		updateZettel(t, repo, z2)
		assert.Equal(t, len(z2.Links), 0, "z2 should not link to z1")
		updateZettel(t, repo, z3)
		assert.Equal(t, len(z3.Links), 0, "z3 should not link to z1")
	})
	t.Run("can remove a zettel and its links", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		zettel := &model.Zettel{
			ID:    "1",
			Title: "Testing Zettel",
		}

		createZettel(t, repo, zettel)

		err := repo.Remove(context.Background(), zettel)
		require.Equal(t, err, nil, "failed to remove zettel")

		err = repo.Get(context.Background(), zettel)
		assert.Equal(t, err, sql.ErrNoRows, "zettel should not exist")
	})
	t.Run("can remove zettel in bulk", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		zettels := []*model.Zettel{
			{
				ID:    "1",
				Title: "Testing Zettel",
			},
			{
				ID:    "2",
				Title: "Testing Zettel",
			},
			{
				ID:    "3",
				Title: "Testing Zettel",
			},
		}

		err := repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettel")

		err = repo.RemoveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to remove zettel in bulk")

		err = repo.Get(context.Background(), zettels[0])
		assert.Equal(t, err, sql.ErrNoRows, "zettel should not exist")

		err = repo.Get(context.Background(), zettels[1])
		assert.Equal(t, err, sql.ErrNoRows, "zettel should not exist")

		err = repo.Get(context.Background(), zettels[2])
		assert.Equal(t, err, sql.ErrNoRows, "zettel should not exist")
	})
}

func TestZettelRepository_List(t *testing.T) {
	t.Run("can list all fleets", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		err := repo.Reset(context.Background())
		require.Equal(t, err, nil, "failed to reset database")

		zettels := []*model.Zettel{
			{
				ID:    "1",
				Title: "Testing Zettel",
			},
			{
				ID:    "2",
				Title: "Testing Zettel",
			},
			{
				ID:    "3",
				Title: "Testing Zettel",
			},
		}

		err = repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettel")

		fleets, err := repo.ListFleet(context.Background())
		require.Equal(t, err, nil, "failed to list fleets")

		assert.Equal(t, len(fleets), 3, "should have no fleets")
	})

	t.Run("can list all permanent ", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		err := repo.Reset(context.Background())
		require.Equal(t, err, nil, "failed to reset database")

		zettels := []*model.Zettel{
			{
				ID:    "1",
				Title: "Testing Zettel",
			},
			{
				ID:    "2",
				Title: "Testing Zettel",
			},
			{
				ID:    "3",
				Title: "Testing Zettel",
			},
		}

		err = repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettel")

		// List All Permanent
		//
		zettels[0].Type = "permanent"
		zettels[0].Path = strings.Replace(zettels[0].Path, "fleet", "permanent", 1)

		zettels[1].Type = "permanent"
		zettels[1].Path = strings.Replace(zettels[1].Path, "fleet", "permanent", 1)

		err = repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to save zettel")

		permanents, err := repo.ListPermanent(context.Background())
		require.Equal(t, err, nil, "failed to list permanent")
		assert.Equal(t, len(permanents), 2, "should have no fleets")
	})

	t.Run("can list all zettels", func(t *testing.T) {
		db := sqltest.CreateDatabase(t, cfg)
		repo := NewZettelRepository(db, cfg)

		err := repo.Reset(context.Background())
		require.Equal(t, err, nil, "failed to reset database")

		zettels := []*model.Zettel{
			{
				ID:    "1",
				Title: "Testing Zettel",
			},
			{
				ID:    "2",
				Title: "Testing Zettel",
			},
			{
				ID:    "3",
				Title: "Testing Zettel",
			},
		}

		err = repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to create zettel")

		// List All Permanent
		//
		zettels[0].Type = "permanent"
		zettels[0].Path = strings.Replace(zettels[0].Path, "fleet", "permanent", 1)

		zettels[1].Type = "permanent"
		zettels[1].Path = strings.Replace(zettels[1].Path, "fleet", "permanent", 1)

		err = repo.SaveBulk(context.Background(), zettels...)
		require.Equal(t, err, nil, "failed to save zettel")

		zettels, err = repo.ListAll(context.Background())
		require.Equal(t, err, nil, "failed to list permanent")
		assert.Equal(t, len(zettels), 3, "should have no fleets")
	})
}

func TestZettelRepository_Search(t *testing.T) {
	t.Run("can search by query", func(t *testing.T) {
		// db := sqltest.CreateDatabase(t, cfg)
		// repo := NewZettelRepository(db, cfg)
		//
		// z1 := &model.Zettel{
		// 	ID:      "1",
		// 	Title:   "Testing Zettel",
		// 	Content: "This is a test",
		// }
		// z2 := &model.Zettel{
		// 	ID:      "2",
		// 	Title:   "Testing Zettel 2",
		// 	Content: "A random test",
		// }
		// z3 := &model.Zettel{
		// 	ID:      "3",
		// 	Title:   "Testing Zettel 3",
		// 	Content: "An example test",
		// }
		//
		// createZettel(t, repo, z1)
		// createZettel(t, repo, z2)
		// createZettel(t, repo, z3)
		//
		// zettels, err := repo.Search(context.Background(), "random")
		// require.Equal(t, err, nil, "failed to search zettels")
		// assert.Equal(t, len(zettels), 1, "should find all zettels")
		//
		// zettels, err = repo.Search(context.Background(), "zettel")
		// require.Equal(t, err, nil, "failed to search zettels")
		// assert.Equal(t, len(zettels), 3, "should find all zettels")
	})
}

func createZettel(t *testing.T, repo ZettelRepository, z *model.Zettel) {
	err := repo.Remove(context.Background(), z)
	if !errors.Is(err, ErrZettelNotFound) {
		require.Equal(t, err, nil, "failed to remove zettel")
	}
	err = repo.Save(context.Background(), z)
	require.Equal(t, err, nil, "failed to create zettel")
}

func updateZettel(t *testing.T, repo ZettelRepository, z *model.Zettel) {
	zet := &model.Zettel{
		ID: z.ID,
	}
	err := repo.Get(context.Background(), zet)
	require.Equal(t, err, nil, "failed to get the zettel from db")

	// update the zettel with the db version
	z.Title = zet.Title
	z.Content = zet.Content
	z.Path = zet.Path
	z.Type = zet.Type
	z.CreatedAt = zet.CreatedAt
	z.UpdatedAt = zet.UpdatedAt

	// Auxiliary fields
	z.Links = zet.Links
	z.Lines = zet.Lines
}

func cleanupZettels(t *testing.T, repo ZettelRepository) {
}
