package repository

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/slugify"
)

type ZettelRepository struct {
	DB *database.Database
}

func NewZettelRepository(db *database.Database) *ZettelRepository {
	return &ZettelRepository{DB: db}
}

func (z *ZettelRepository) Get(ctx context.Context, zettel *model.Zettel) error {
	var query string

	if zettel.ID != "" {
		query = `select * from zettel where id = ?`
	} else if zettel.Path != "" {
		query = `select * from zettel where path = ?`
	} else {
		return errors.New("error: zettel id or path must be provided")
	}

	err := z.DB.DB.GetContext(ctx, zettel, query, zettel.ID)
	if err != nil {
		return err
	}

	links := []*model.Zettel{}
	query = `
	SELECT z2.* FROM link l JOIN zettel z2 ON l.link_id = z2.id WHERE l.zettel_id = ?
	`
	err = z.DB.DB.SelectContext(ctx, &links, query, zettel.ID)
	if err != nil {
		return err
	}

	zettel.Links = links

	return nil
}

func (r *ZettelRepository) Create(ctx context.Context, z *model.Zettel) error {
	query := `
  insert into zettel (id, title, content, type, path)
	values (:id, :title, :content, :type, :path)
	on conflict (id) do
	update set title = :title, content = :content, type = :type, path = :path
  returning id, title, content, type, path, created_at, updated_at
  `

	// Set the zettel default values
	if z.Title == "" {
		return errors.New("error: title cannot be empty")
	}
	if z.ID == "" {
		z.ID = isosec()
	}
	if z.Path == "" {
		z.Path = config.FLEET_PATH + "/" + slugify.Slug(z.Title) + "." + z.ID + ".md"
	}
	if z.Content == "" {
		z.Content = emptyContent(z.Title)
	}
	if z.Type == "" {
		z.Type = "fleet"
	}

	rows, err := r.DB.DB.NamedQueryContext(ctx, query, z)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(z); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (r *ZettelRepository) CreateBulk(ctx context.Context, zettels ...*model.Zettel) error {
	query := `
  insert into zettel (id, title, content, type, path)
	values (:id, :title, :content, :type, :path)
	on conflict (id) do
	update set title = :title, content = :content, type = :type, path = :path
  `

	// Set the zettel default values
	for _, z := range zettels {
		if z.Title == "" {
			return errors.New("error: title cannot be empty")
		}
		if z.ID == "" {
			z.ID = isosec()
		}
		if z.Path == "" {
			z.Path = config.FLEET_PATH + "/" + slugify.Slug(z.Title) + "." + z.ID + ".md"
		}
		if z.Content == "" {
			z.Content = emptyContent(z.Title)
		}
		if z.Type == "" {
			z.Type = "fleet"
		}
	}

	_, err := r.DB.DB.NamedExecContext(ctx, query, zettels)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZettelRepository) Link(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
	query := `
	insert into link (zettel_id, link_id) values (:zettel_id, :link_id)
	on conflict (zettel_id, link_id) do nothing
	`

	links := make([]model.Link, len(zettels))
	for i, z2 := range zettels {
		links[i] = model.Link{
			From: z1.ID,
			To:   z2.ID,
		}
	}
	_, err := z.DB.DB.NamedExecContext(ctx, query, links)
	if err != nil {
		return err
	}

	dbLinks := []*model.Zettel{}
	err = z.DB.DB.SelectContext(ctx, &dbLinks, `select z2.* from link l join zettel z2 on l.link_id = z2.id where l.zettel_id = ?`, z1.ID)
	if err != nil {
		return err
	}

	z1.Links = dbLinks

	return nil
}

func (z *ZettelRepository) LinkBulk(ctx context.Context, links ...*model.Link) error {
	query := `
	insert into link (zettel_id, link_id) values (:zettel_id, :link_id)
	on conflict (zettel_id, link_id) do nothing
	`

	_, err := z.DB.DB.NamedExecContext(ctx, query, links)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZettelRepository) Unlink(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
	query := `
	delete from link
	where zettel_id = ? and link_id in (?)
	`

	// Convert slice of Zettel into a comma separated string of IDs.
	ids := make([]string, len(zettels))
	for i, zet := range zettels {
		ids[i] = zet.ID
	}

	// Replace ? with the actual list of IDs.
	query, args, err := sqlx.In(query, z1.ID, ids)
	if err != nil {
		return err
	}

	// sqlx.In returns queries with ? bindvars, we can rebind it for our
	// database.
	query = z.DB.DB.Rebind(query)

	_, err = z.DB.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	dbLinks := []*model.Zettel{}
	err = z.DB.DB.SelectContext(ctx, &dbLinks, `select z2.* from link l join zettel z2 on l.link_id = z2.id where l.zettel_id = ?`, z1.ID)
	if err != nil {
		return err
	}

	z1.Links = dbLinks

	return nil
}

// emptyContent returns an empty content for a zettel, which has the following
// structure:
// # <title>
// <empty line>
// <empty line>
func emptyContent(title string) string {
	return "# " + title + "\n\n\n"
}

// isosec generates now timestamps like 20220605165935(0-99999) using the atomic
// package to generate id's, avoiding collisions
func isosec() string {
	var counter uint64
	return fmt.Sprintf("%s%01d", time.Now().Format("20060102150405"), atomic.AddUint64(&counter, 1)%100000)
}
