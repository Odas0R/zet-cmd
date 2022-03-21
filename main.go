package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gosimple/slug"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "zet",
		Version: "v0.1",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "odas0r",
				Email: "guilherme.odas0r@gmail.com",
			},
		},
		Usage:     "A zettelkasten under a terminal approach",
		UsageText: "//TODO",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return nil
					}

					title := strings.Join(c.Args().Slice(), " ")
					id := time.Now().Unix()
					slug := slug.Make(title)

					fileName := fmt.Sprintf("%s.%d.md", slug, id)

					filePath := fmt.Sprintf("%s/%s", fleetPath, fileName)

					zettel, err := Zettel{
						ID:       id,
						Title:    title,
						FileName: fileName,
						Path:     filePath,
						Tags:     []string{},
						Links:    []string{},
					}.Create()

					if err != nil {
						return err
					}

					// opens the zettel on the specified $EDITOR
					zettel.Open()

					return nil
				},
			},
			{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "",
				Action: func(c *cli.Context) error {

					query, err := filepath.Abs("./github.com/zet-cmd/scripts/query")
					if err != nil {
						return err
					}

					fmt.Println(query)

					// Execute the script query
					cmd := exec.Command(query)
					// cmd.Start()

					bytes, err := cmd.Output()
					if err != nil {
						return err
					}

					filePath := string(bytes[:])

					fmt.Println(filePath)

					return nil
				},
			},
			{
				Name:    "backlog",
				Aliases: []string{"bg"},
				Usage:   "",
				Action: func(c *cli.Context) error {

					files, err := ioutil.ReadDir("/tmp/")
					if err != nil {
						log.Fatal(err)
					}

					for _, file := range files {
						fmt.Println(file.Name(), file.IsDir())
					}

					// https://github.com/AlecAivazis/survey
					color := ""
					prompt := &survey.Select{
						Message: "Choose a color:",
						Options: []string{"red", "blue", "green"},
					}

					survey.AskOne(prompt, &color)

					return nil
				},
			},
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "",
				Action: func(c *cli.Context) error {
					config := &Config{}

					err := config.Init()
					if err != nil {
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
