package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/urfave/cli/v2"
)

const (
	fleetPath     = "./example/fleet"
	permanentPath = "./example/fleet"
	templatesPath = "./example/templates"
	historyPath   = "./example/templates"

  // scripts
  open = "./scripts/open"
  query = "./scripts/query"
)

type Zettel struct {
	ID       int64
	FileName string
	Title    string
	Tags     []string
	Links    []string
}

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

					fmt.Println(fileName)

					zettel := &Zettel{
						ID:       id,
						Title:    title,
						FileName: fileName,
						Tags:     []string{},
						Links:    []string{},
					}

					// parse the template
					tmpl, err := template.ParseFiles(fmt.Sprintf("%s/zet.tmpl.md", templatesPath))
					if err != nil {
						return err
					}

					filePath := fmt.Sprintf("%s/%s", fleetPath, fileName)

					// create the zettel file
					f, err := os.Create(filePath)
					if err != nil {
						return err
					}

					// put the given title to the zettel
					err = tmpl.Execute(f, zettel)
					if err != nil {
						return err
					}
					f.Close()

					// open the respective zettel
          cmd := exec.Command(open, filePath)
          cmd.Start()

					return nil
				},
			},
			{
				Name:    "query",
				Aliases: []string{"q"},
				Usage:   "",
				Action: func(c *cli.Context) error {

					// Execute the script query
          cmd := exec.Command(query)
          cmd.Start()

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
