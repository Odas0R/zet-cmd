package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/samber/lo"
)

type Zettel struct {
	ID       int64
	Title    string
	Type     string
	FileName string
	Slug     string
	Path     string
	Tags     []string
	Links    []string
	Lines    []string
}

func (z *Zettel) New(c *Config) error {
	if z.Title == "" {
		return errors.New("title cannot be empty")
	}

	// assign the zettel metadata
	if z.ID == 0 {
		z.ID = time.Now().Local().UTC().Unix()
	}

	z.Slug = slug.Make(z.Title)
	z.FileName = fmt.Sprintf("%s.%d.md", z.Slug, z.ID)
	z.Path = fmt.Sprintf("%s/%s", c.Sub.Fleet, z.FileName)
	z.Type = "fleet" // can be "fleet" or "permanent"
	z.Tags = []string{}
	z.Links = []string{}

	// create the zettel file
	file, err := os.Create(z.Path)
	if err != nil {
		return err
	}

	// parse the template file of the zettel
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/zet.tmpl.md", c.Sub.Templates))
	if err != nil {
		return err
	}

	// put the given title to the zettel
	err = tmpl.Execute(file, z)

	if err != nil {
		return err
	}
	file.Close()

	// Set the lines of the file
	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	z.Lines = lines

	return nil
}

func (z *Zettel) Read(c *Config) error {
	if z.Path == "" {
		return errors.New("error: zettel path cannot be empty")
	}

	if !strings.Contains(z.Path, c.Root) {
		return errors.New("error: file is not under the root path")
	}

	if fileExists := FileExists(z.Path); !fileExists {
		return errors.New("error: zettel does not exist on given path")
	}

	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	z.Lines = lines

	// convert ID to int64
	basename := filepath.Base(z.Path)
	indexOne := strings.Index(basename, ".")
	indexTwo := strings.LastIndex(basename, ".")
	if indexOne == indexTwo {
		return errors.New("error: invalid zettel id")
	}

	idStr := basename[indexOne+1 : indexTwo]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	z.ID = id
	z.Title = strings.ReplaceAll(z.Lines[0], "# ", "")
	z.Slug = strings.Split(filepath.Base(z.Path), ".")[0]
	z.FileName = filepath.Base(z.Path)

	foldersNames := strings.Split(z.Path, "/")
	z.Type = foldersNames[len(foldersNames)-2]

	var isOnLinksSection = false
	z.Links = lo.Reduce(z.Lines, func(acc []string, line string, i int) []string {
		if i == len(z.Lines)-1 {
			return acc
		}

		if line == "## Links" {
			isOnLinksSection = true
		}

		if isOnLinksSection {
			indexOne := strings.LastIndex(line, "(") + 1
			indexTwo := strings.LastIndex(line, ")")

			if indexOne != -1 && indexTwo != -1 {
				link := line[indexOne:indexTwo]
				return append(acc, link)
			}
		}

		return acc
	}, []string{})

	z.Tags = lo.Filter(strings.Split(z.Lines[len(z.Lines)-1], " "), func(tag string, _ int) bool {
		m := regexp.MustCompile(`^#\w+$`)
		return m.FindString(tag) == tag
	})

	return nil
}

func (z *Zettel) Link(zettel *Zettel) error {
	// read the file and check if there's a link
	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	link := fmt.Sprintf("- [%s](%s)", zettel.Title, zettel.Path)

	// check if link already exists
	if hasLink := lo.Contains(lines, link); hasLink {
		return errors.New("error: cannot have duplicated links")
	}

	// append on the file
	lineToInsert := lo.IndexOf(z.Lines, "## Links") + 1
	link = fmt.Sprintf("\n- [%s](%s)", zettel.Title, zettel.Path)
	if err := AppendLine(z.Path, link, lineToInsert); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Repair(c *Config, h *History) error {
	//
	// Get Metadata
	//
	basename := filepath.Base(z.Path)
	indexOne := strings.Index(basename, ".")
	indexTwo := strings.LastIndex(basename, ".")
	if indexOne == indexTwo {
		return errors.New("error: invalid zettel id")
	}

	idStr := basename[indexOne+1 : indexTwo]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	z.ID = id
	z.Title = strings.ReplaceAll(z.Lines[0], "# ", "")
	z.Slug = slug.Make(strings.ReplaceAll(z.Title, "_", ""))
	z.FileName = fmt.Sprintf("%s.%d.md", z.Slug, z.ID)

	foldersNames := strings.Split(z.Path, "/")
	z.Type = foldersNames[len(foldersNames)-2]

	//
	// Fix title
	//

	oldPath := z.Path
	newPath := fmt.Sprintf("%s/%s/%s", c.Root, z.Type, z.FileName)

  // return if oldPath is equal to newPath, just means that the title wasn't
  // changed
  if oldPath == newPath {
    return nil
  }
  
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	z.Path = newPath

	// Fix the history
	h.Delete(oldPath)
	h.Insert(newPath)

	//
	// Fix the broken links on other zettels
	//
	cmd := exec.Command("/bin/bash", c.Scripts.FindLinks, idStr, c.Sub.Fleet, c.Sub.Permanent)

	data, err := cmd.Output()
	if err != nil {
		return err
	}

	newLink := fmt.Sprintf("- [%s](%s)", z.Title, z.Path)
	entries := strings.Split(bytes.NewBuffer(data).String(), "\n")
	entries = entries[:len(entries)-1] // remove last element

	// go through every entry and update the dirty links
	for _, entry := range entries {
		values := strings.Split(entry, ":")

		lineNr, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}

		lineIndex := lineNr - 1
		filepath := values[1]

		zettel := &Zettel{Path: filepath}
		if err := zettel.Read(c); err != nil {
			return err
		}

		zettel.Lines = lo.ReplaceAll(zettel.Lines, zettel.Lines[lineIndex], newLink)

		if err := zettel.Write(); err != nil {
			return err
		}
	}

	// update links
	var isOnLinksSection = false
	z.Links = lo.Reduce(z.Lines, func(acc []string, line string, i int) []string {
		if i == len(z.Lines)-1 {
			return acc
		}

		if line == "## Links" {
			isOnLinksSection = true
		}

		if isOnLinksSection {
			indexOne := strings.LastIndex(line, "(") + 1
			indexTwo := strings.LastIndex(line, ")")

			if indexOne != -1 && indexTwo != -1 {
				link := line[indexOne:indexTwo]
				return append(acc, link)
			}
		}

		return acc
	}, []string{})

  // update tags
	z.Tags = lo.Filter(strings.Split(z.Lines[len(z.Lines)-1], " "), func(tag string, _ int) bool {
		m := regexp.MustCompile(`^#\w+$`)
		return m.FindString(tag) == tag
	})

	return nil
}

func (z *Zettel) ReadLines() error {
	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	// update lines
	z.Lines = lines

	return nil
}

func (z *Zettel) Write() error {
	output := strings.Join(z.Lines, "\n")

	if err := ioutil.WriteFile(z.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

// func (z *Zettel) Delete() error {
// 	return nil
// }

func (z *Zettel) WriteLine(index int, newLine string) error {
	// modify z.Lines
	copy(z.Lines[index:], []string{newLine})

	output := strings.Join(z.Lines, "\n")

	if err := ioutil.WriteFile(z.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Permanent(c *Config, h *History) error {
	newPath := fmt.Sprintf("%s/%s", c.Sub.Permanent, z.FileName)
	if err := os.Rename(z.Path, newPath); err != nil {
		return err
	}

	// update the path
	z.Path = newPath

	// fix all broken links
	if err := z.Repair(c, h); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Open(c *Config, lineNr int) error {
	if err := Open(c, z.Path, lineNr); err != nil {
		return err
	}

	return nil
}
