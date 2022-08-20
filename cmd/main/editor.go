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

func Query(initial string) (string, int, error) {
	cmd := exec.Command("/bin/bash", config.Scripts.Query, initial, config.Sub.Fleet, config.Sub.Permanent)
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

func Open(path string, lineNr int) error {
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
