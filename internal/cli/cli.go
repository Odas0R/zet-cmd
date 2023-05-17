package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/muxit-studio/color"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/repository"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/fs"
	"github.com/odas0r/zet/pkg/slugify"
	"github.com/urfave/cli/v2"
)

func New(db *database.Database) *cli.App {

	app := &cli.App{
		Name:    "zet",
		Version: "0.1",
		Authors: []*cli.Author{
			{
				Name:  "odas0r",
				Email: "guilherme.odas0r@gmail.com",
			},
		},
		Usage:                "A zettelkasten under a terminal approach",
		UsageText:            "A simple way to manage your zettelkasten using neovim (telescope) and fzf",
		Flags:                []cli.Flag{},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "Create a new zettel",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					zr := repository.NewZettelRepository(db)

					newZet := &model.Zettel{
						Title: strings.Join(c.Args().Slice(), " "),
					}

					if err := zr.Create(context.Background(), newZet); err != nil {
						log.Fatalf("error: failed to create zettel: %v", err)
					}

					if err := fs.WriteToFile(newZet.Path, newZet.Content); err != nil {
						log.Fatalf("error: failed to write to file: %v", err)
					}

					// place on stdout the path of the new zettel
					io.WriteString(c.App.Writer, newZet.Path)

					return nil
				},
			},
			{
				Name:  "query",
				Usage: "Init telescope to search for zettels",
				Action: func(c *cli.Context) error {
					query := strings.Join(c.Args().Slice(), " ")

					var cmd string

					if query == "" {
						if fs.HasNvimSession() {
							cmd = "nvim --server \"$NVIM_SOCKET\" --remote-send \":ZetQuery<CR>\""
						} else {
							cmd = "nvim --listen \"$NVIM_SOCKET\" -c \":ZetQuery\""
						}
					} else {
						if fs.HasNvimSession() {
							cmd = "nvim --server \"$NVIM_SOCKET\" --remote-send \":ZetQuery " + query + "<CR>\""
						} else {
							cmd = "nvim --listen \"$NVIM_SOCKET\" -c \":ZetQuery " + query + "\""
						}
					}

					if err := fs.Exec(cmd); err != nil {
						log.Fatalf("error: failed to execute command: %v", err)
					}

					return nil
				},
			},
			{
				Name:  "save",
				Usage: "Inserts or updates the given zettel to the database",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}
					path := c.Args().Slice()[0]

					zr := repository.NewZettelRepository(db)
					zet := &model.Zettel{
						Path: path,
					}

					if !zet.IsValid() {
						log.Fatalf("error: zettel on path %s is not valid", path)
					}

					query := `select * from zettel where path = ?`

					if err := db.DB.GetContext(context.Background(), zet, query, path); err != nil {
						if !errors.Is(err, sql.ErrNoRows) {
							return err
						}
					}

					if err := zet.Read(); err != nil {
						log.Fatalf("error: failed to read from file: %v", err)
					}

					if err := zr.Create(context.Background(), zet); err != nil {
						log.Fatalf("error: failed to create zettel: %v", err)
					}

					// Add links if there are any
					if len(zet.Links) > 0 {
						if err := zr.Link(context.Background(), zet, zet.Links); err != nil {
							log.Fatalf("error: failed to link zettel: %v", err)
						}
					}

					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "Sync the filesystem with the database and does some fixing on the side",
				Action: func(c *cli.Context) error {
					fleet := fs.List(config.FLEET_PATH)
					perm := fs.List(config.PERMANENT_PATH)

					paths := append(fleet, perm...)

					var zettels []*model.Zettel
					var counter uint64
					for _, path := range paths {
						// 1. Initialize the new zettel and read
						zet := &model.Zettel{
							Path: path,
						}
						if err := zet.Read(); err != nil {
							log.Fatalf("error: failed to read from file: %v", err)
						}

						if zet.ID == "" {
							// generate a new ID for the zettel, using the current timestamp
							// and a counter to avoid collisions
							zet.ID = fmt.Sprintf("%s%01d", time.Now().Format("20060102150405"), atomic.AddUint64(&counter, 1)%100000)

							// create a new path for the zettel using the new ID
							basename := slugify.Slug(zet.Title) + "." + zet.ID + ".md"
							zet.Path = filepath.Join(filepath.Dir(zet.Path), basename)

							if err := fs.WriteToFile(zet.Path, zet.Content); err != nil {
								log.Fatalf("error: failed to write to file: %v", err)
							}

							// remove the old file
							if err := fs.Remove(path); err != nil {
								log.Fatalf("error: failed to remove old zettel: %v", err)
							}

							fmt.Printf(color.BGreen("[NEW_PATH]: %s\n"), zet.Path)
						}

						zettels = append(zettels, zet)
					}

					zr := repository.NewZettelRepository(db)

					fmt.Println(color.BGreen("[LOG]: INDEXING THE ZETTELKASTEN..."))

					// 5. Do a CreateBulk for all the zettels
					// 6. Do a LinkBulk for all the zettels

					return nil
				},
			},
		},
	}

	return app
}
