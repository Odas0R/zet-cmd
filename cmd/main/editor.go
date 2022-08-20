package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Query(initial string, c *Config) (string, int, error) {
	cmd := exec.Command("/bin/bash", c.Scripts.Query, initial, c.Sub.Fleet, c.Sub.Permanent)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	value, err := cmd.Output()
	if err != nil {
		return "", 0, err
	}

	str := strings.Split(bytes.NewBuffer(value).String(), ":")
	lineNr, err := strconv.Atoi(str[0])
	if err != nil {
		return "", 0, err
	}

	path := str[1]

	return strings.TrimSpace(path), lineNr, nil
}

func Ripgrep(query string, c *Config) ([]string, error) {
	cmd := exec.Command("/bin/bash", c.Scripts.Ripgrep, strings.TrimSpace(query), c.Sub.Fleet, c.Sub.Permanent)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	value, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(bytes.NewBuffer(value).String(), "\n")

	return lines, nil
}

func Open(config *Config, path string, lineNr int) error {
	if fileExists := FileExists(path); !fileExists {
		return errors.New("error: file does not exist")
	}

	cmd := exec.Command("/bin/bash", config.Scripts.Open, path, fmt.Sprintf("%v", lineNr))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
