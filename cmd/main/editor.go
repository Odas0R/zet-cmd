package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

func Query(c *Config, tags []string) error {
	// Format tags
	tagsFormatted := lo.Map(tags, func(tag string, i int) string {
		return fmt.Sprintf("\"%s\"", tag)
	})
	tagsStr := strings.Join(tagsFormatted, ",")

	// Execute the script query
	cmd := exec.Command("/bin/bash", c.Scripts.Query, c.Root, tagsStr)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Open(config *Config, path string) error {
	if fileExists := FileExists(path); !fileExists {
		return errors.New("error: file does not exist")
	}

	cmd := exec.Command("/bin/bash", config.Scripts.Open, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
