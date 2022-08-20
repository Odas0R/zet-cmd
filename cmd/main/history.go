package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
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

	if err := h.Read(); err != nil {
		return err
	}

	return nil
}

func (h *History) Query() (string, error) {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return "", err
	}

	// transform into titles
	var titles = make([]string, 0, len(lines))
	var zettels = make([]*Zettel, 0, len(lines))

	for _, line := range lines {
		zet := &Zettel{Path: line}
		if err := zet.Read(); err != nil {
			return "", err
		}

		titles = append(titles, fmt.Sprintf("%s | %s", color.UYellow(zet.Lines[0]), strings.Join(zet.Tags, " ")))
		zettels = append(zettels, zet)
	}

	cmd := exec.Command("/bin/bash", config.Scripts.Fzf, columnize.SimpleFormat(titles), "70%")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputStr := strings.TrimSpace(bytes.NewBuffer(output).String())
	zettelPath, ok := lo.Find(zettels, func(zet *Zettel) bool {
		return zet.Lines[0] == outputStr
	})

	if !ok {
		return "", nil
	}

	return zettelPath.Path, nil
}

func (h *History) Insert(path string) error {
	if path == "" {
		return errors.New("error: path cannot be empty")
	}

	// Validate zettel
	zettel := &Zettel{Path: path}
	if err := zettel.Read(); err != nil {
		return err
	}

	lines, err := ReadLines(h.Path)
	if err != nil {
		return err
	}

	// append the new file
	lines = lo.Filter(lines, func(line string, i int) bool {
		return line != path
	})
	lines = append(lines, path)

	// write to the history
	output := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(h.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func (h *History) Delete(path string) error {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return err
	}

	lines = lo.Filter(lines, func(line string, i int) bool {
		return line != path
	})

	output := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(h.Path, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func (h *History) Read() error {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return err
	}

	// update lines
	h.Lines = lines

	return nil
}

func (h *History) Open() error {
	if err := Open(h.Path, 0); err != nil {
		return err
	}

	return nil
}
