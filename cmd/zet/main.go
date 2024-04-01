package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/internal/model"
	"github.com/odas0r/zet/internal/repository"
	"github.com/odas0r/zet/pkg/database"
	"github.com/odas0r/zet/pkg/fs"
	"github.com/urfave/cli/v2"
)

const (
	rootDir     = "/home/odas0r/github.com/odas0r/zet"
	databaseUrl = "file:/home/odas0r/github.com/odas0r/zet-cmd/zettel.db"
	// rootDir     = "/tmp/zet"
	// databaseUrl = "file:/tmp/zet/zettel.db"
)

func main() {
	db := database.NewDatabase(database.NewDatabaseOptions{
		URL:                databaseUrl,
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	config := config.New(rootDir)
	zr := repository.NewZettelRepository(db, config)

	app := &cli.App{
		Name:    "zet",
		Version: "0.1",
		Authors: []*cli.Author{
			{
				Name:  "odas0r",
				Email: "guilherme@muxit.co",
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "raw",
						Usage: "Create a new zettel and output the path to stdout",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					title := strings.Join(c.Args().Slice(), " ")

					zet, err := New(zr, title)
					if err != nil {
						log.Fatalf("error: failed to create new zettel: %v", err)
					}

					if c.Bool("raw") {
						io.WriteString(os.Stdout, zet.Path)
						return nil
					}

					yes := fs.InputConfirm("Do you want to open the zettel?")
					if yes {
						if err := fs.Editor(zet.Path); err != nil {
							log.Fatalf("error: failed to open file: %v", err)
						}
					}

					return nil
				},
			},
			{
				Name:  "open",
				Usage: "Opens the zettel by the given path",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					path := c.Args().Slice()[0]

					zet := &model.Zettel{
						Path: path,
					}

					if !zet.IsValid(config) {
						log.Fatalf("error: invalid zettel with given path %s", path)
					}

					if err := fs.Editor(zet.Path); err != nil {
						log.Fatalf("error: failed to open file: %v", err)
					}

					return nil
				},
			},
			{
				Name:  "search",
				Usage: "Search for zettels using sqlite3 fs5 extension",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					query := strings.Join(c.Args().Slice(), " ")

					zettels, err := Search(zr, query)
					if err != nil {
						log.Fatalf("error: failed to search for zettels: %v", err)
					}

					for _, zet := range zettels {
						fmt.Fprintf(c.App.Writer, "%s\n", zet.Path)
					}

					return nil
				},
			},
			{
				Name: "remove",
				Aliases: []string{
					"rm",
				},
				Usage: "Removes the given zettel from the database and from the filesystem",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}
					path := c.Args().Slice()[0]

					if err := Remove(zr, path); err != nil {
						log.Fatalf("error: failed to remove zettel: %v", err)
					}

					return nil
				},
			},
			{
				Name:  "history",
				Usage: "Retrieves the last 50 opened zettel",
				Action: func(_ *cli.Context) error {
					zettels, err := History(zr)
					if err != nil {
						log.Fatalf("error: failed to query the history: %v", err)
					}

					for _, zettel := range zettels {
						io.WriteString(os.Stdout, zettel.Path+"\n")
					}

					return nil
				},
			},
			{
				Name:  "backlog",
				Usage: "Retrieves all the fleet of zettels",
				Action: func(_ *cli.Context) error {
					zettels, err := Backlog(zr)
					if err != nil {
						log.Fatalf("error: failed to query the backlog: %v", err)
					}

					for _, zettel := range zettels {
						io.WriteString(os.Stdout, zettel.Path+"\n")
					}

					return nil
				},
			},
			{
				Name:  "links",
				Usage: "Retrieves all the links of a zettel",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}
					path := c.Args().Slice()[0]

					zettels, err := Links(zr, path)
					if err != nil {
						return err
					}

					for _, zettel := range zettels {
						io.WriteString(os.Stdout, zettel.Path+"\n")
					}

					return nil
				},
			},
			{
				Name:  "backlinks",
				Usage: "Retrieves all the backlinks of a zettel",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}
					path := c.Args().Slice()[0]

					zettels, err := BackLinks(zr, path)
					if err != nil {
						return err
					}

					for _, zettel := range zettels {
						io.WriteString(os.Stdout, zettel.Path+"\n")
					}

					return nil
				},
			},
			{
				Name:  "brokenlinks",
				Usage: "Retrieves all the brokenlinks of a zettel",
				Action: func(_ *cli.Context) error {
					zettels, err := BrokenLinks(zr)
					if err != nil {
						log.Fatalf("error: failed to query all the brokenlinks of a zettel: %v", err)
					}

					for _, zettel := range zettels {
						io.WriteString(os.Stdout, zettel.Path+"\n")
					}

					return nil
				},
			},
			{
				Name:  "last",
				Usage: "Retrieves the last opened zettel",
				Action: func(_ *cli.Context) error {
					// fetch the last edited zettel
					zet, err := Last(zr)
					if err != nil {
						log.Fatalf("error: failed to query the last opened zettel: %v", err)
					}

					io.WriteString(os.Stdout, zet.Path)

					return nil
				},
			},
			{
				Name:  "save",
				Usage: "Inserts or updates the given zettel to the database, and some repairs",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}
					path := c.Args().Slice()[0]

					zet, err := Save(zr, path)
					if err != nil {
						log.Fatalf("error: failed to save zettel: %v", err)
					}

					io.WriteString(os.Stdout, zet.Path)

					return nil
				},
			},
			{
				// indexing phase
				Name:  "sync",
				Usage: "Sync the filesystem with the database and does some fixing on the side",
				Action: func(_ *cli.Context) error {
					if err := Sync(zr); err != nil {
						log.Fatalf("error: failed to sync zettels: %v", err)
					}

					fmt.Println("Synced! :)")

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
