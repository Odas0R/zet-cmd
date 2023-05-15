package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {

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

					title := strings.Join(c.Args().Slice(), " ")

					fmt.Println(title)

					return nil
				},
			},
		},
	}

	return app
}
