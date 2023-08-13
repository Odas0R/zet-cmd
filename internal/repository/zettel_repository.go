package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
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

	// Save works for both fleet and permanent, if you to make a zettel permanent
	// you just need to update the path :)
	Save(ctx context.Context, zettel *model.Zettel) error
	SaveBulk(ctx context.Context, zettels ...*model.Zettel) error
	Link(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	LinkBulk(ctx context.Context, links ...*model.Link) error
	Unlink(ctx context.Context, zettel *model.Zettel, links []*model.Zettel) error
	Remove(ctx context.Context, zettel *model.Zettel) error
	RemoveBulk(ctx context.Context, zettels ...*model.Zettel) error
	LastOpened(ctx context.Context, zettel *model.Zettel) error
	InsertHistory(ctx context.Context, zettel *model.Zettel) error
	History(ctx context.Context) ([]*model.Zettel, error)
	ListFleet(ctx context.Context) ([]*model.Zettel, error)
	ListPermanent(ctx context.Context) ([]*model.Zettel, error)
	ListAll(ctx context.Context) ([]*model.Zettel, error)
	Backlinks(ctx context.Context, zet *model.Zettel) ([]*model.Zettel, error)
	Search(ctx context.Context, query string) ([]*model.Zettel, error)
	Reset(ctx context.Context) error
	Config() *config.Config
}

type zettelRepository struct {
	config *config.Config
	DB     *database.Database
}

func NewZettelRepository(db *database.Database, config *config.Config) ZettelRepository {
	return &zettelRepository{
		config: config,
		DB:     db,
	}
}

func (zr *zettelRepository) Config() *config.Config {
	return zr.config
}

func (zr *zettelRepository) Get(ctx context.Context, zettel *model.Zettel) error {
	var query string

	if zettel.ID != "" {
		query = `select * from zettel where id = ?`
	} else if zettel.Path != "" {
		query = `select * from zettel where path = ?`
	} else {
		return errors.New("error: zettel id or path must be provided")
	}

	err := zr.DB.DB.GetContext(ctx, zettel, query, zettel.ID)
	if err != nil {
		return err
	}

	links := []*model.Zettel{}
	query = `
	SELECT z2.* FROM link l JOIN zettel z2 ON l.link_id = z2.id WHERE l.zettel_id = ?
	`
	err = zr.DB.DB.SelectContext(ctx, &links, query, zettel.ID)
	if err != nil {
		return err
	}

	zettel.Links = links
	zettel.Lines = strings.Split(zettel.Content, "\n")

	return nil
}

func (zr *zettelRepository) Save(ctx context.Context, z *model.Zettel) error {
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
		z.Path = zr.config.FleetRoot + "/" + slug.Make(z.Title) + "." + z.ID + ".md"
	}
	if z.Content == "" {
		z.Content = emptyContent(z.Title)
		z.Lines = strings.Split(z.Content, "\n")
	}
	if z.Type == "" {
		z.Type = "fleet"
	}

	rows, err := zr.DB.DB.NamedQueryContext(ctx, query, z)
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

func (zr *zettelRepository) SaveBulk(ctx context.Context, zettels ...*model.Zettel) error {
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
			z.Path = zr.config.FleetRoot + "/" + slug.Make(z.Title) + "." + z.ID + ".md"
		}
		if z.Content == "" {
			z.Content = emptyContent(z.Title)
		}
		if z.Type == "" {
			z.Type = "fleet"
		}
	}

	_, err := zr.DB.DB.NamedExecContext(ctx, query, zettels)
	if err != nil {
		return err
	}

	return nil
}

func (zr *zettelRepository) Link(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
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
	_, err := zr.DB.DB.NamedExecContext(ctx, query, links)
	if err != nil {
		return err
	}

	dbLinks := []*model.Zettel{}
	err = zr.DB.DB.SelectContext(ctx, &dbLinks, `select z2.* from link l join zettel z2 on l.link_id = z2.id where l.zettel_id = ?`, z1.ID)
	if err != nil {
		return err
	}

	z1.Links = dbLinks

	return nil
}

func (zr *zettelRepository) LinkBulk(ctx context.Context, links ...*model.Link) error {
	query := `
	insert into link (zettel_id, link_id) values (:zettel_id, :link_id)
	on conflict (zettel_id, link_id) do nothing
	`

	_, err := zr.DB.DB.NamedExecContext(ctx, query, links)
	if err != nil {
		return err
	}

	return nil
}

func (zr *zettelRepository) Unlink(ctx context.Context, z1 *model.Zettel, zettels []*model.Zettel) error {
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
	query = zr.DB.DB.Rebind(query)

	_, err = zr.DB.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	dbLinks := []*model.Zettel{}
	err = zr.DB.DB.SelectContext(ctx, &dbLinks, `select z2.* from link l join zettel z2 on l.link_id = z2.id where l.zettel_id = ?`, z1.ID)
	if err != nil {
		return err
	}

	z1.Links = dbLinks

	return nil
}

func (zr *zettelRepository) Remove(ctx context.Context, zettel *model.Zettel) error {
	query := `delete from zettel where id = :id`

	res, err := zr.DB.DB.NamedExecContext(ctx, query, zettel)
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

func (zr *zettelRepository) RemoveBulk(ctx context.Context, zettels ...*model.Zettel) error {
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
	query = zr.DB.DB.Rebind(query)

	_, err = zr.DB.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (zr *zettelRepository) LastOpened(ctx context.Context, zettel *model.Zettel) error {
	query := `
	select z.* from history as h
	inner join zettel as z on h.zettel_id = z.id
	order by h.updated_at desc limit 1
	`

	err := zr.DB.DB.GetContext(ctx, zettel, query)
	if err != nil {
		return err
	}

	return nil
}

func (zr *zettelRepository) InsertHistory(ctx context.Context, zet *model.Zettel) error {
	query := `
  insert into history (zettel_id)
  values (?) on conflict (zettel_id) do update
  set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ')
  `

	_, err := zr.DB.DB.ExecContext(ctx, query, zet.ID)
	if err != nil {
		return err
	}

	return nil
}

func (zr *zettelRepository) History(ctx context.Context) ([]*model.Zettel, error) {
	query := `
  select z.* from zettel as z
  inner join history as h on z.id = h.zettel_id
  order by h.updated_at desc limit 50
  `

	zettels := []*model.Zettel{}
	err := zr.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (zr *zettelRepository) ListFleet(ctx context.Context) ([]*model.Zettel, error) {
	query := `select * from zettel where type = 'fleet' order by updated_at desc`

	zettels := []*model.Zettel{}
	err := zr.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (zr *zettelRepository) ListPermanent(ctx context.Context) ([]*model.Zettel, error) {
	query := `select * from zettel where type = 'permanent' order by updated_at desc`

	zettels := []*model.Zettel{}
	err := zr.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (zr *zettelRepository) ListAll(ctx context.Context) ([]*model.Zettel, error) {
	query := `select * from zettel where type = 'fleet' or type = 'permanent' order by updated_at desc`

	zettels := []*model.Zettel{}
	err := zr.DB.DB.SelectContext(ctx, &zettels, query)
	if err != nil {
		return nil, err
	}

	return zettels, nil
}

func (zr *zettelRepository) Reset(ctx context.Context) error {
	query := `delete from zettel returning *`
	_, err := zr.DB.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (zr *zettelRepository) Backlinks(ctx context.Context, zet *model.Zettel) ([]*model.Zettel, error) {
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

	rows, err := zr.DB.DB.NamedQueryContext(ctx, query, zet)
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

func (zr *zettelRepository) Search(ctx context.Context, query string) ([]*model.Zettel, error) {
	q := `
	select
		z.id,
		z.title,
		z.content,
		z.path,
		z.created_at,
	  z.updated_at
	from zettel z
	  join zettel_fts zf on (zf.rowid = z.id)
	where zettel_fts match ?
	order by rank
	`

	var zettels []*model.Zettel
	err := zr.DB.DB.SelectContext(ctx, &zettels, q, query)
	if err != nil {
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
	// if is in test mode, add counter to avoid collisions
	if os.Getenv("TEST") == "true" {
		return fmt.Sprintf("%s%01d", time.Now().Format("20060102150405"), atomic.AddUint64(&counter, 1)%100000)
	}

	return time.Now().Format("20060102150405")
}
