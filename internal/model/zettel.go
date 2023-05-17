package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/pkg/fs"
	"github.com/odas0r/zet/pkg/slugify"
)

type Zettel struct {
	ID        string
	Title     string
	Content   string
	Lines     []string
	Path      string
	Type      string
	Links     []*Zettel
	CreatedAt Time `db:"created_at"`
	UpdatedAt Time `db:"updated_at"`
}

// IsValid checks if file is a zettel and if it exists.
func (z *Zettel) IsValid() bool {
	// Verify if the zettel is valid by checking its ID
	if z.ID != "" && z.Path == "" {
		permZettels := fs.List(config.PERMANENT_PATH)
		fleetZettels := fs.List(config.FLEET_PATH)

		zettels := append(permZettels, fleetZettels...)

		for _, zettel := range zettels {
			id := "." + z.ID + ".md"
			if strings.Contains(zettel, id) {
				return true
			}
		}

		return false
	}

	// Verify if the zettel is valid by checking its path
	return fs.Exists(z.Path) && (strings.Contains(z.Path, config.FLEET_PATH) ||
		strings.Contains(z.Path, config.PERMANENT_PATH))
}

// Read reads a zettel from the disk and gets all the metadata from it. Useful
// to query data from the file and insert into a database.
func (z *Zettel) Read() error {
	if !z.IsValid() {
		return fmt.Errorf("error: zettel is not valid")
	}

	lines, err := fs.ReadLines(z.Path)
	if err != nil {
		return err
	}

	z.ID = z.readId()
	z.Title = strings.TrimPrefix(lines[0], "# ")
	z.Content = strings.Join(lines, "\n")
	z.Lines = lines
	z.Type = z.readType()

	// read links, get all titles from [[wikilinks]] with the format [[#
	// <title>|<id>]]
	var links []*Zettel
	for _, line := range lines {
		results := fs.MatchAllSubstrings("[[", "]]", line)
		for _, result := range results {
			if result != "" {
				link := &Zettel{
					ID: result,
				}
				if link.IsValid() {
					links = append(links, link)
				}
			}
		}
	}

	z.Links = links

	return nil
}

func (z *Zettel) Write() error {
	return fs.Write(z.Path, z.Content)
}

func (z *Zettel) Repair() error {
	oldPath := z.Path

	z.Title = strings.TrimPrefix(z.Lines[0], "# ")
	z.Path = filepath.Join(filepath.Dir(z.Path), slugify.Slug(z.Title)+"."+z.ID+".md")

	if oldPath != z.Path {
		if err := z.Write(); err != nil {
			return err
		}

		if err := fs.Remove(oldPath); err != nil {
			return err
		}
	}

	return nil
}

func (z *Zettel) IsEqual(z2 *Zettel) bool {
	return z.ID == z2.ID && z.Title == z2.Title && z.Content == z2.Content && z.Path == z2.Path && z.Type == z2.Type
}

func (z *Zettel) readId() string {
	fileName := filepath.Base(z.Path)
	return fs.MatchSubstring(".", ".", fileName)
}

func (z *Zettel) readType() string {
	var typ string
	if strings.Contains(z.Path, config.FLEET_PATH) {
		typ = "fleet"
	} else if strings.Contains(z.Path, config.PERMANENT_PATH) {
		typ = "permanent"
	}
	return typ
}
