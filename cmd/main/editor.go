package main

import (
	"errors"
	"os"
	"os/exec"
)

func Query(initial string, c *Config) error {
	cmd := exec.Command("/bin/bash", c.Scripts.Query, initial, c.Sub.Fleet, c.Sub.Permanent)
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
