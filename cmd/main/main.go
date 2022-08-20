package main

import (
	"errors"
	"log"
	"os"
	"strings"

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
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Create a new zettel",
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
				Name:    "last",
				Aliases: []string{"l"},
				Usage:   "",
				Action: func(c *cli.Context) error {
					if err := history.Read(); err != nil {
						return err
					}

					zettel := &Zettel{Path: history.Lines[len(history.Lines)-1]}
					if err := zettel.Read(); err != nil {
						return err
					}

					if err := zettel.Open(0); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "history",
				Aliases: []string{"h"},
				Usage:   "",
				Action: func(c *cli.Context) error {
					path, err := history.Query()
					if err != nil {
						return err
					}

					zettel := &Zettel{Path: path}

					// validate zettel
					if err := zettel.Read(); err != nil {
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

							if err := history.Insert(path); err != nil {
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

							if err := history.Delete(path); err != nil {
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
				Name:    "backlog",
				Aliases: []string{"bg"},
				Usage:   "",
				Action: func(c *cli.Context) error {
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

					if err := zettel.Read(); err != nil {
						return err
					}

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
