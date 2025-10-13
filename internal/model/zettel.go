package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/odas0r/zet/internal/config"
	"github.com/odas0r/zet/pkg/fs"
)

type Zettel struct {
	ID        string `db:"id" json:"id"`
	Slug      string `db:"slug" json:"slug"`
	Title     string `db:"title" json:"title"`
	Path      string `db:"path" json:"path"`
	CreatedAt Time   `db:"created_at" json:"createdAt"`
	UpdatedAt Time   `db:"updated_at" json:"updatedAt"`

	// prevent content from beign shared on json responses
	Content string `db:"content" json:"-"`

	// Auxiliary fields (not stored in the database)
	Lines []string  `json:"-"`
	Links []*Zettel `json:"-"`
}

// IsValid checks if file is a zettel and if it exists.
func (z *Zettel) IsValid(cfg *config.Config) bool {
	if z.ID != "" {
		zettels := fs.List(cfg.Path)
		for _, zettel := range zettels {
			if strings.Contains(zettel, z.ID) {
				return true
			}
		}
	}
	return fs.Exists(z.Path)
}

// Read reads a zettel from the disk and gets all the metadata from it. Useful
// to query data from the file and insert into a database.
func (z *Zettel) Read(cfg *config.Config) error {
	if !z.IsValid(cfg) {
		return fmt.Errorf("error: zettel is not valid")
	}

	lines, err := fs.ReadLines(z.Path)
	if err != nil {
		return err
	}

	z.ID = z.readId()
	z.Title = strings.TrimPrefix(lines[0], "# ")
	z.Slug = slug.Make(z.Title)
	z.Content = strings.Join(lines, "\n")
	z.Lines = lines

	// read links, get all slugs from [[wikilinks]] like [[slug-link]]
	var links []*Zettel
	var mapLinks = make(map[string]bool)
	for _, line := range lines {
		results := fs.MatchAllSubstrings("[[", "]]", line)
		for _, result := range results {
			if _, ok := mapLinks[result]; !ok && result != z.Slug && result != "" {
				var s = result
				if !slug.IsSlug(s) {
					s = slug.Make(s)
				}
				link := &Zettel{
					Slug: s,
				}
				links = append(links, link)
				mapLinks[result] = true
			}
		}
	}

	z.Links = links

	return nil
}

func (z *Zettel) Write() error {
	return fs.Write(z.Path, z.Content)
}

func (z *Zettel) WriteLine(line string) error {
	z.Content += "\n" + line
	return z.Write()
}

func (z *Zettel) IsEqual(z2 *Zettel) bool {
	return z.ID == z2.ID && z.Title == z2.Title && z.Content == z2.Content && z.Path == z2.Path
}

func (z *Zettel) readId() string {
	fileName := filepath.Base(z.Path)
	// Remove the extension .md
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
