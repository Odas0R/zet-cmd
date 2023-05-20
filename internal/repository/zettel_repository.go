package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/pkg/database"
)

var counter uint64

// Errors
var (
	ErrZettelNotFound = errors.New("error: zettel not found")
)

type ZettelRepository interface {
	Get(ctx context.Context, zettel *model.Zettel) error
	Create(ctx context.Context, zettel *model.Zettel) error
	CreateBulk(ctx context.Context, zettels ...*model.Zettel) error
	Link(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	LinkBulk(ctx context.Context, links ...*model.Link) error
	Unlink(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	Remove(ctx context.Context, zettel *model.Zettel) error
	RemoveBulk(ctx context.Context, zettels ...*model.Zettel) error
	LastOpened(ctx context.Context, zettel *model.Zettel) error
	History(ctx context.Context) ([]*model.Zettel, error)
	ListFleet(ctx context.Context) ([]*model.Zettel, error)
	Backlinks(ctx context.Context, zet *model.Zettel) ([]*model.Zettel, error)
}

type zettelRepository struct {
	DB *database.Database
}

func NewZettelRepository(db *database.Database) ZettelRepository {
	return &zettelRepository{DB: db}
}

func (z *zettelRepository) Get(ctx context.Context, zettel *model.Zettel) error {
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
	zettel.Lines = strings.Split(zettel.Content, "\n")

	return nil
}

func (r *zettelRepository) Create(ctx context.Context, z *model.Zettel) error {
	query := `
  insert into zettel (id, title, content, type, path)
	values (:id, :title, :content, :type, :path)
	on conflict (id) do
	update set title = excluded.title, content = excluded.content, type = excluded.type, path = excluded.path
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
		z.Path = config.FLEET_PATH + "/" + slug.Make(z.Title) + "." + z.ID + ".md"
	}
	if z.Content == "" {
		z.Content = emptyContent(z.Title)
		z.Lines = strings.Split(z.Content, "\n")
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

func (r *zettelRepository) CreateBulk(ctx context.Context, zettels ...*model.Zettel) error {
	query := `
  insert into zettel (id, title, content, type, path)
	values (:id, :title, :content, :type, :path)
	on conflict(id) do update set
	title = excluded.title,
	content = excluded.content,
	type = excluded.type,
	path = excluded.path
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
			z.Path = config.FLEET_PATH + "/" + slug.Make(z.Title) + "." + z.ID + ".md"
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

func (z *zettelRepository) Link(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
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

func (z *zettelRepository) LinkBulk(ctx context.Context, links ...*model.Link) error {
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

func (z *zettelRepository) Unlink(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
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

func (z *zettelRepository) Remove(ctx context.Context, zettel *model.Zettel) error {
	query := `delete from zettel where id = :id`

	res, err := z.DB.DB.NamedExecContext(ctx, query, zettel)
	if err != nil {
		return err
	}

	nr, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if nr == 0 {
		return ErrZettelNotFound
	}

	return nil
}

func (z *zettelRepository) RemoveBulk(ctx context.Context, zettels ...*model.Zettel) error {
	query := `delete from zettel where id in (?)`

	// Convert slice of Zettel into a comma separated string of IDs.
	ids := make([]string, len(zettels))
	for i, zet := range zettels {
		ids[i] = zet.ID
	}

	// Replace ? with the actual list of IDs.
	query, args, err := sqlx.In(query, ids)
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

	return nil
}

func (z *zettelRepository) LastOpened(ctx context.Context, zettel *model.Zettel) error {
	query := `select * from zettel order by updated_at desc limit 1`

	err := z.DB.DB.GetContext(ctx, zettel, query)
	if err != nil {
		return err
	}

	return nil
}

func (z *zettelRepository) History(ctx context.Context) ([]*model.Zettel, error) {
	query := `select * from zettel order by updated_at desc limit 50`

	zettels := []*model.Zettel{}
	err := z.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (z *zettelRepository) ListFleet(ctx context.Context) ([]*model.Zettel, error) {
	query := `select * from zettel where type = 'fleet' order by updated_at desc`

	zettels := []*model.Zettel{}
	err := z.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (z *zettelRepository) Backlinks(ctx context.Context, zet *model.Zettel) ([]*model.Zettel, error) {
	query := `
	select z.*
	from zettel z
	join link l on z.id = l.zettel_id
	where l.link_id = (select id from zettel where path = :path);
	`

	if zet.Path == "" {
		query = `
		select z.*
		from zettel z
		join link l on z.id = l.zettel_id
		where l.link_id = :id;
		`
	}

	rows, err := z.DB.DB.NamedQueryContext(ctx, query, zet)
	if err != nil {
		return nil, err
	}

	zettels := []*model.Zettel{}
	for rows.Next() {
		zet := &model.Zettel{}
		err := rows.StructScan(zet)
		if err != nil {
			return nil, err
		}

		zettels = append(zettels, zet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return zettels, nil
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
	return fmt.Sprintf("%s%01d", time.Now().Format("20060102150405"), atomic.AddUint64(&counter, 1)%100000)
}
