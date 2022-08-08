package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

const FILE_NAME = ".history"

type History struct {
	Path string
	Root string
}

func (h *History) Init() error {
	if h.Root == "" {
		return errors.New("error: history cannot have Root empty")
	}

	h.Path = fmt.Sprintf("%s/%s", h.Root, FILE_NAME)

	if err := CreateFile("", h.Path); err != nil {
		return err
	}

	return nil
}

func (h *History) Query(config *Config) (string, error) {
	lines, err := ReadLines(h.Path)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("/bin/bash", config.Scripts.Fzf, strings.Join(lines, "\n"))

	value, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return bytes.NewBuffer(value).String(), nil
}

func (h *History) Insert(path string) error {
	if path == "" {
		return errors.New("error: path cannot be empty")
	}

	if !strings.Contains(path, h.Root) {
		return errors.New("error: path must be valid")
	}

	if fileExists := FileExists(path); !fileExists {
		return errors.New("error: file does not exist")
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

func (h *History) Open(c *Config) error {
	if err := Open(c, h.Path); err != nil {
		return err
	}

	return nil
}
