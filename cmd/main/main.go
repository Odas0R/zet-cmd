package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/odas0r/zet/cmd/color"
	"github.com/odas0r/zet/cmd/columnize"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func main() {

	// initialize config
	if err := config.Init(os.Getenv("ZET")); err != nil {
		log.Fatalf("error: failed to initialize config %V", err)
	}

	// initialize history
	if err := history.Init(os.Getenv("ZET"), ".history"); err != nil {
		log.Fatalf("error: failed to initialize history %V", err)
	}

	app := &cli.App{
		Name:    "zet",
		Version: "v0.1",
		Authors: []*cli.Author{
			&cli.Author{
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

					title := strings.Join(c.Args().Slice(), " ")

					zettel := &Zettel{Title: title}
					if err := zettel.New(); err != nil {
						return err
					}

					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},

			{
				Name:  "link",
				Usage: "Link two or more zettels",
				Action: func(c *cli.Context) error {
					zettel, err := history.Query()
					if err != nil {
						return err
					}

					zettels, err := history.QueryMany()
					if err != nil {
						return err
					}

					for _, zettelToBeLinked := range zettels {
						if err := zettel.Link(zettelToBeLinked); err != nil {
							return err
						}
					}

          // Open the zettel that the links were written
					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "",
				Action: func(c *cli.Context) error {
					path, line, err := Query("")
					if err != nil {
						return err
					}

					zettel := &Zettel{Path: path}
					if err := zettel.Open(line); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "backlog",
				Usage: "Query the zettelkasten backlog/inbox and open a specific fleet note",
				Action: func(c *cli.Context) error {
					files, err := ioutil.ReadDir(config.Sub.Fleet)
					if err != nil {
						log.Fatal(err)
					}

					// sort files by access time
					sort.Slice(files, func(i, j int) bool {
						return files[i].ModTime().After(files[j].ModTime())
					})

					var rows = make([]string, 0, len(files))
					var zettels = make([]*Zettel, 0, len(files))
					for _, file := range files {
						zettel := &Zettel{Path: fmt.Sprintf("%s/%s", config.Sub.Fleet, file.Name())}

						err := zettel.ReadLines()
						if err != nil {
							continue
						}

						// TODO: might wanna "columnize" e.g.  fmt.Sprintf(%s | %s, col1,
						// col2)
						row := color.UYellow(zettel.Lines[0])

						rows = append(rows, row)
						zettels = append(zettels, zettel)
					}

					output, err := Fzf(columnize.SimpleFormat(rows), "70%", "Backlog > ")
					if err != nil {
						return err
					}

					zettel, ok := lo.Find(zettels, func(zet *Zettel) bool {
						return strings.HasPrefix(output, zet.Lines[0])
					})
					if !ok {
						return errors.New("error: no zettel found")
					}

					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "last",
				Usage: "",
				Action: func(c *cli.Context) error {
					if err := history.Read(); err != nil {
						return err
					}

					zettel := &Zettel{Path: history.Lines[0]}
					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "history",
				Usage: "",
				Action: func(c *cli.Context) error {
					zettel, err := history.Query()
					if err != nil {
						return err
					}

					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "insert",
						Usage: "",
						Action: func(c *cli.Context) error {
							if c.NArg() == 0 {
								return errors.New("error: empty arguments")
							}

							path := strings.Join(c.Args().Slice(), "")

							zettel := &Zettel{Path: path}
							if err := history.Insert(zettel); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "delete",
						Usage: "",
						Action: func(c *cli.Context) error {
							if c.NArg() == 0 {
								return errors.New("error: empty arguments")
							}

							path := strings.Join(c.Args().Slice(), "")

							zettel := &Zettel{Path: path}
							if err := history.Delete(zettel); err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:  "edit",
						Usage: "",
						Action: func(c *cli.Context) error {
							if err := history.Open(); err != nil {
								return err
							}

							return nil
						},
					},
				},
			},
			{
				Name:    "delete",
				Usage:   "",
				Aliases: []string{"rm"},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return errors.New("error: empty arguments")
					}

					path := strings.Join(c.Args().Slice(), "")
					zettel := &Zettel{Path: path}

					if err := zettel.Delete(); err != nil {
						return err
					}

					// clear editor buffer
					if err := DeleteBuffer(); err != nil {
						return err
					}

					//
					// Open query after deletion
					//

					path, line, err := Query("")
					if err != nil {
						return err
					}

					zettel = &Zettel{Path: path}
					if err := zettel.Open(line); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "repair",
				Usage: "",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return errors.New("error: empty arguments")
					}

					path := strings.Join(c.Args().Slice(), "")
					zettel := &Zettel{Path: path}

					if err := zettel.Repair(); err != nil {
						return err
					}

					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
