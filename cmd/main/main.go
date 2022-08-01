package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {

	config := &Config{Path: "/tmp/foo"}

	// initialize config
	config.Init()

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
					if err := zettel.New(config); err != nil {
						return err
					}

					if err := zettel.Open(); err != nil {
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
					if err := Query(); err != nil {
						return err
					}

					return nil
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
