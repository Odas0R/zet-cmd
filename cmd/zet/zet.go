package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/repository"
	"github.com/odas0r/zet/pkg/fs"
)

func New(zr repository.ZettelRepository, title string) (*model.Zettel, error) {
	zet := &model.Zettel{
		Title: title,
	}

	if err := zr.Save(context.Background(), zet); err != nil {
		return nil, err
	}

	if err := fs.Write(zet.Path, zet.Content); err != nil {
		return nil, err
	}

	return zet, nil
}

func Search(zr repository.ZettelRepository, query string) ([]*model.Zettel, error) {
	return zr.Search(context.Background(), query)
}

func Remove(zr repository.ZettelRepository, path string) error {
	zet := &model.Zettel{
		Path: path,
	}

	if err := zr.Remove(context.Background(), zet); err != nil {
		return err
	}

	if fs.Exists(zet.Path) {
		if err := fs.Remove(zet.Path); err != nil {
			return err
		}
	}

	return nil
}

func History(zr repository.ZettelRepository) ([]*model.Zettel, error) {
	return zr.History(context.Background())
}

func Backlog(zr repository.ZettelRepository) ([]*model.Zettel, error) {
	return zr.ListFleet(context.Background())
}

func Links(zr repository.ZettelRepository, path string) ([]*model.Zettel, error) {
	zet := &model.Zettel{
		Path: path,
	}

	err := zr.Get(context.Background(), zet)
	if err != nil {
		return nil, err
	}

	return zet.Links, nil
}

func BackLinks(zr repository.ZettelRepository, path string) ([]*model.Zettel, error) {
	zet := &model.Zettel{
		Path: path,
	}

	zettels, err := zr.Backlinks(context.Background(), zet)
	if err != nil {
		return nil, err
	}
	return zettels, nil
}

// How it works?
//
// - A broken link is when [[<empty>]] or [[<invalid_id>]]
//
// TODO: improve the return so that it shows on which line and column
// the broken link is so that you can use with quickfix list
func BrokenLinks(zr repository.ZettelRepository) ([]*model.Zettel, error) {
	zettels, err := zr.ListAll(context.Background())
	if err != nil {
		return nil, err
	}

	var brokenZettels []*model.Zettel

	for _, zet := range zettels {
		if err := zet.Read(zr.Config()); err != nil {
			return nil, err
		}

		ok, err := zet.HasBrokenLinks(zr.Config())
		if err != nil {
			return nil, err
		}

		if ok {
			brokenZettels = append(brokenZettels, zet)
		}
	}
	return brokenZettels, nil
}

func Last(zr repository.ZettelRepository) (*model.Zettel, error) {
	zet := &model.Zettel{}

	if err := zr.LastOpened(context.Background(), zet); err != nil {
		return nil, err
	}

	return zet, nil
}

func Save(zr repository.ZettelRepository, path string) (*model.Zettel, error) {
	zet := &model.Zettel{Path: path}

	// Get all the zettel metadata
	if err := zet.Read(zr.Config()); err != nil {
		return nil, err
	}

	if zet.ID == "" {
		zet.ID = time.Now().Format("20060102150405")
	}

	// Do repairs if needed
	if err := zet.Repair(); err != nil {
		return nil, err
	}

	// Update the database with the new zettel, based on the ID
	if err := zr.Save(context.Background(), zet); err != nil {
		return nil, err
	}

	if err := zr.InsertHistory(context.Background(), zet); err != nil {
		return nil, err
	}

	// Add links if there are any
	if len(zet.Links) > 0 {
		if err := zr.Link(context.Background(), zet, zet.Links); err != nil {
			return nil, err
		}
	}

	return zet, nil
}

func Sync(zr repository.ZettelRepository) error {
	cfg := zr.Config()

	fleet := fs.List(cfg.FleetRoot)
	perm := fs.List(cfg.PermanentRoot)

	paths := append(fleet, perm...)

	var zettels []*model.Zettel
	var counter uint64

	for _, path := range paths {
		// 1. Initialize the new zettel and read
		zet := &model.Zettel{
			Path: path,
		}
		if err := zet.Read(cfg); err != nil {
			return err
		}

		if zet.ID == "" {
			// generate a new ID for the zettel, using the current timestamp
			// and a counter to avoid collisions
			zet.ID = fmt.Sprintf("%s%01d", time.Now().Format("20060102150405"), atomic.AddUint64(&counter, 1)%100000)
		}

		if err := zet.Repair(); err != nil {
			return err
		}

		zettels = append(zettels, zet)
	}

	if err := zr.SaveBulk(context.Background(), zettels...); err != nil {
		return err
	}

	var links []*model.Link
	for _, zet := range zettels {
		for _, zetLink := range zet.Links {
			link := &model.Link{
				From: zet.ID,
				To:   zetLink.ID,
			}
			links = append(links, link)
		}
	}

	if len(links) > 0 {
		if err := zr.LinkBulk(context.Background(), links...); err != nil {
			return err
		}
	}

	// cleaning up phase
	//

	dbZettel, err := zr.ListAll(context.Background())
	if err != nil {
		return err
	}

	var toRemove []*model.Zettel
	for _, zet := range dbZettel {
		if !fs.Exists(zet.Path) {
			toRemove = append(toRemove, zet)
		}
	}

	if len(toRemove) > 0 {
		if err := zr.RemoveBulk(context.Background(), toRemove...); err != nil {
			return err
		}
	}

	return nil
}
