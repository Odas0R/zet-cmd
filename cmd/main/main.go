package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {

	config := &Config{Root: os.Getenv("ZET")}
	history := &History{Root: os.Getenv("ZET")}

	// initialize config
	if err := config.Init(); err != nil {
		log.Fatalf("error: failed to initialize config %V", err)
	}

	// initialize history
	if err := history.Init(); err != nil {
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
					if err := zettel.New(config); err != nil {
						return err
					}

					if err := zettel.Open(config, 0); err != nil {
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

					// query := strings.Join(c.Args().Slice(), " ")
					//
					//      lines, err := Ripgrep(query, config)
					//      if err != nil {
					//        return err
					//      }
					//
					//      fmt.Printf("lines: %v\n", lines)

					path, line, err := Query("", config)
					if err != nil {
						return err
					}

					zettel := &Zettel{Path: path}
					if err := zettel.Open(config, line); err != nil {
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
					path, err := history.Query(config)
					if err != nil {
						return err
					}

					zettel := &Zettel{Path: path}

					// validate zettel
					if err := zettel.Read(config); err != nil {
						return err
					}

					if err := zettel.Open(config, 0); err != nil {
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
							if err := history.Open(config); err != nil {
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

					if err := zettel.Read(config); err != nil {
						return err
					}

					if err := zettel.Repair(config, history); err != nil {
						return err
					}

					if err := zettel.Open(config, 0); err != nil {
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
