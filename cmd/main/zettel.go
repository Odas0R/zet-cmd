package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/samber/lo"
)

type Zettel struct {
	// Metadata that is fetched from the `body` of the file and `path`
	ID       int64
	Title    string
	Type     string
	FileName string
	Slug     string
	Path     string

	// Represents the lines of the file of the `zettel`
	Lines []string

	// Represents all associated `zettels`
	Links []*Zettel
}

func (z *Zettel) New() error {
	if z.Title == "" {
		return errors.New("title cannot be empty")
	}

	// assign the zettel metadata
	if z.ID == 0 {
		z.ID = time.Now().Local().UTC().Unix()
	}

	z.Slug = slug.Make(z.Title)
	z.FileName = fmt.Sprintf("%s.%d.md", z.Slug, z.ID)
	z.Path = fmt.Sprintf("%s/%s", config.Sub.Fleet, z.FileName)
	z.Type = "fleet" // can be "fleet" or "permanent"
	z.Links = []*Zettel{}

	// create the zettel file
	file, err := os.Create(z.Path)
	if err != nil {
		return err
	}

	// parse the template file of the zettel
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/zet.tmpl.md", config.Sub.Templates))
	if err != nil {
		return err
	}

	// put the given title to the zettel
	if err := tmpl.Execute(file, z); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	// Set the lines of the file
	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	z.Lines = lines

	return nil
}

func (z *Zettel) Read() error {
	if err := z.ReadMetadata(); err != nil {
		return err
	}

	if err := z.ReadLinks(); err != nil {
		return err
	}
	return nil
}

func (z *Zettel) ReadMetadata() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	if err := z.ReadLines(); err != nil {
		return err
	}

	fileName := filepath.Base(z.Path)
	idStr, ok := MatchSubstring(".", ".", fileName)
	if !ok {
		return errors.New("error: zettel has an invalid id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	z.ID = id
	z.Title = strings.Replace(z.Lines[0], "# ", "", 1)
	z.Slug = slug.Make(z.Title)

	foldersNames := strings.Split(z.Path, "/")
	z.Type = foldersNames[len(foldersNames)-2]
	z.FileName = fmt.Sprintf("%s.%d.md", z.Slug, z.ID)

	return nil
}

func (z *Zettel) ReadLinks() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	if err := z.ReadLines(); err != nil {
		return err
	}

	zettels, err := z.getLinks()
	if err != nil {
		return err
	}

	z.Links = zettels

	return nil
}

func (z *Zettel) getLinks() ([]*Zettel, error) {
	if err := z.ReadLines(); err != nil {
		return nil, err
	}

	_, linksIndex, ok := lo.FindIndexOf(z.Lines, func(line string) bool {
		return strings.Contains(line, "Links")
	})

	if !ok {
		return nil, fmt.Errorf("error: given zettel has no '## Links' section")
	}

	if linksIndex == -1 {
		return nil, errors.New("error: there is no ## Links section on given zettel")
	}

	zettels := []*Zettel{}
	for index, line := range z.Lines {
		if index > linksIndex {
			path, ok := ValidateLinkPath(line)
			if ok {
				zet := &Zettel{Path: path}
				if zet.IsValid() {

					zettels = append(zettels, zet)
				}
			}
		}
	}
	return zettels, nil
}

func (z *Zettel) Link(zettel *Zettel) error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}
	if !zettel.IsValid() {
		return errors.New("error: given zettel is invalid")
	}

	// Get both zettels data
	if err := z.ReadLines(); err != nil {
		return err
	}
	if err := zettel.ReadMetadata(); err != nil {
		return err
	}

	if z.Path == zettel.Path {
		return errors.New("error: cannot link the same file")
	}

	// check if link already exists
	linkToInsert := fmt.Sprintf("- [%s](%s)", zettel.Title, zettel.Path)
	if lo.Contains(z.Lines, linkToInsert) {
		return errors.New("error: cannot have duplicated links")
	}

	_, linksIndex, _ := lo.FindIndexOf(z.Lines, func(line string) bool {
		return strings.Contains(line, "Links")
	})

	links, err := z.getLinks()
	if err != nil {
		return err
	}

	linksStringArr := lo.Map(links, func(z *Zettel, _ int) string {
		z.ReadMetadata()
		return fmt.Sprintf("- [%s](%s)", z.Title, z.Path)
	})

	//
	// Format links and insert the new link
	//

	lines := z.Lines[:linksIndex+1]
	lines = append(lines, "") // add a <new-line>
	lines = append(lines, linksStringArr...)
	lines = append(lines, linkToInsert) // add link to the end of the file
	lines = append(lines, "")           // add a <new-line>
	lines = append(lines, "")           // add a <new-line>

	z.Lines = lines

	if err := z.Write(); err != nil {
		return err
	}

	if err := z.ReadLinks(); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Repair() (error, bool) {
  titleHasChanged := false

	if !z.IsValid() {
		return errors.New("error: invalid zettel"), titleHasChanged
	}

	// Zemove the old zettel from history
	if err := history.Delete(z); err != nil {
		return err, titleHasChanged
	}

	// Get new metadata
	if err := z.ReadMetadata(); err != nil {
		return err, titleHasChanged
	}

	// Title was changed, update the path
	if filepath.Base(z.Path) != z.FileName {
		oldPath := z.Path
		z.Path = fmt.Sprintf("%s/%s/%s", config.Path, z.Type, z.FileName)
		if err := os.Rename(oldPath, z.Path); err != nil {
			return err, false
		}
    titleHasChanged = true
	}

	// Fix the history
	if err := history.Insert(z); err != nil {
		return err, titleHasChanged
	}

	results, ok := Grep(strconv.FormatInt(z.ID, 10))
	if !ok {
		return nil, titleHasChanged
	}

	//
	// Fix all dity links
	//
	for _, result := range results {
		zettel := &Zettel{Path: result.Path}
		ok := zettel.IsValid()
		if !ok {
			return errors.New("error: file path on link is not a valid zettel"), false
		}

		if err := zettel.ReadLines(); err != nil {
			return err, titleHasChanged
		}

		index := result.LineNr - 1
		zettel.Lines[index] = fmt.Sprintf("- [%s](%s)", z.Title, z.Path)

		if err := zettel.Write(); err != nil {
			return err, titleHasChanged
		}
	}

	return nil, titleHasChanged
}

func (z *Zettel) ReadLines() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	lines, err := ReadLines(z.Path)
	if err != nil {
		return err
	}

	// update lines
	z.Lines = lines

	return nil
}

func (z *Zettel) Write() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	output := strings.Join(z.Lines, "\n")

	if err := ioutil.WriteFile(z.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Delete() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	if err := z.ReadMetadata(); err != nil {
		return err
	}

	// Remove deleted zettel from history, if exists
	if err := history.Delete(z); err != nil {
		return err
	}

	idStr := strconv.FormatInt(z.ID, 10)
	results, ok := Grep(idStr)
	if !ok {
		// Delete file
		if err := os.Remove(z.Path); err != nil {
			return err
		}

		return nil
	}

	//
	// Fix all dirty links
	//
	for _, result := range results {
		zettel := &Zettel{Path: result.Path}
		ok := zettel.IsValid()
		if !ok {
			return errors.New("error: file path on link is not a valid zettel")
		}

		if err := zettel.ReadLines(); err != nil {
			return err
		}

		index := result.LineNr - 1
		zettel.Lines = append(zettel.Lines[:index], zettel.Lines[index+1:]...)

		if err := zettel.Write(); err != nil {
			return err
		}
	}

	// Delete file
	if err := os.Remove(z.Path); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) WriteLine(index int, newLine string) error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	if err := z.ReadLines(); err != nil {
		return err
	}

	// modify z.Lines
	copy(z.Lines[index:], []string{newLine})

	if err := z.Write(); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Permanent() error {
	if !z.IsValid() {
		return errors.New("error: invalid zettel")
	}

	// update metadata
	if err := z.ReadMetadata(); err != nil {
		return nil
	}

	// Move zettel into the the /permanent directory
	oldPath := z.Path
	z.Path = fmt.Sprintf("%s/%s", config.Sub.Permanent, z.FileName)
	if err := os.Rename(oldPath, z.Path); err != nil {
		return err
	}

	// fix all broken links
	if err, _ := z.Repair(); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) Open(lineNr int) error {
	if !z.IsValid() {
		return errors.New("error: zettel is not valid")
	}

	if err := Edit(z.Path, lineNr); err != nil {
		return err
	}

	return nil
}

func (z *Zettel) IsValid() bool {
	return strings.Contains(z.Path, config.Path) && FileExists(z.Path)
}

// Validates if a string is formatted accordingly, and if the string is a valid
// link, by validating the path, checking if it's a valid zettel.
//
// - [$title]($path)
func ValidateLinkPath(str string) (string, bool) {
	path, ok := MatchSubstring("(", ")", str)
	if !ok {
		return "", false
	}

	zettel := &Zettel{Path: path}
	if !zettel.IsValid() {
		return "", false
	}

	return path, true
}
