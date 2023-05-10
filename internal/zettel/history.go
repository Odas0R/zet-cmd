package zettel

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/odas0r/zet/cmd/color"
	"github.com/odas0r/zet/cmd/columnize"
	"github.com/samber/lo"
)

var history = &History{}

type History struct {
	Path  string
	Lines []string
}

func (h *History) Init(root string, fileName string) error {
	if root == "" {
		return errors.New("error: history cannot have root empty")
	}
	if fileName == "" {
		return errors.New("error: history cannot have filename empty")
	}

	h.Path = fmt.Sprintf("%s/%s", root, fileName)

	if err := NewFile("", h.Path); err != nil {
		return err
	}

	if err := h.ReadLines(); err != nil {
		return err
	}

	return nil
}

func (h *History) Query() (*Zettel, error) {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return &Zettel{}, err
	}

	// transform into titles
	var titles = make([]string, 0, len(lines))
	var zettels = make([]*Zettel, 0, len(lines))

	for _, line := range lines {
		zet := &Zettel{Path: line}

		if err := zet.ReadLines(); err != nil {
			return nil, err
		}

		// TODO: might wanna "columnize" e.g.  fmt.Sprintf(%s | %s, col1, col2)
		titles = append(titles, color.UYellow(zet.Lines[0]))
		zettels = append(zettels, zet)
	}

	output, err := Fzf(columnize.SimpleFormat(titles), "70%", "History > ")
	if err != nil {
		return &Zettel{}, err
	}

	zettel, ok := lo.Find(zettels, func(zet *Zettel) bool {
		return strings.HasPrefix(output, zet.Lines[0])
	})
	if !ok {
		return &Zettel{}, nil
	}

	return zettel, nil
}

func (h *History) QueryMany() ([]*Zettel, error) {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return []*Zettel{}, err
	}

	// transform into titles
	var titles = make([]string, 0, len(lines))
	var zettels = make([]*Zettel, 0, len(lines))

	for _, line := range lines {
		zet := &Zettel{Path: line}

		if err := zet.ReadLines(); err != nil {
			return nil, err
		}

		// Show the minutes ago the file was opened!!

		// TODO: might wanna "columnize" e.g.  fmt.Sprintf(%s | %s, col1, col2)
		titles = append(titles, color.UYellow(zet.Lines[0]))
		zettels = append(zettels, zet)
	}

	output, err := FzfMultipleSelection(columnize.SimpleFormat(titles), "70%", "History > ")
	if err != nil {
		return []*Zettel{}, err
	}

	zettels = lo.Filter(zettels, func(zet *Zettel, _ int) bool {
		_, ok := lo.Find(output, func(path string) bool {
			return strings.HasPrefix(path, zet.Lines[0])
		})
		return ok
	})

	return zettels, nil
}

func (h *History) Insert(zettel *Zettel) error {
	if !zettel.IsValid() {
		return errors.New("error: given zettel was not valid")
	}

	if err := h.ReadLines(); err != nil {
		return err
	}

	if len(h.Lines) == 50 {
		return errors.New("error: history cannot have more than 50 zettels")
	}

	if err := history.Delete(zettel); err != nil {
		return err
	}

	h.Lines = append([]string{zettel.Path}, h.Lines...)

	if err := h.Write(); err != nil {
		return err
	}

	return nil
}

func (h *History) Delete(zettel *Zettel) error {
	if !zettel.IsValid() {
		return errors.New("error: zettel invalid")
	}

	lines, err := ReadLines(h.Path)
	if err != nil {
		return err
	}

	lines = lo.Filter(lines, func(line string, i int) bool {
		return line != zettel.Path
	})

	h.Lines = lines

	if err := h.Write(); err != nil {
		return err
	}

	return nil
}

func (h *History) ReadLines() error {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return err
	}

	// update lines
	h.Lines = lines

	return nil
}

func (h *History) Write() error {
	output := strings.Join(h.Lines, "\n")

	if err := ioutil.WriteFile(h.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func (h *History) Clear() error {
	h.Lines = []string{}

	if err := h.Write(); err != nil {
		return err
	}

	return nil
}

func (h *History) Open() error {
	if err := Edit(h.Path, 0); err != nil {
		return err
	}

	return nil
}
