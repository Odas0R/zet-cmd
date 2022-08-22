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

func Fzf(data string, layout string, prompt string) (string, error) {
	cmd := exec.Command("/bin/bash", config.Scripts.Fzf, data, layout, prompt)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(bytes.NewBuffer(output).String()), nil
}

func FzfMultipleSelection(data string, layout string, prompt string) ([]string, error) {
	cmd := exec.Command("/bin/bash", config.Scripts.FzfMulti, data, layout, prompt)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	outputArr := strings.Split(strings.TrimSpace(bytes.NewBuffer(output).String()), "\n")

	return outputArr, nil
}

func GrepLinksById(id int64) ([]string, error) {
	cmd := exec.Command("/bin/bash", config.Scripts.FindLinks, fmt.Sprintf("%d", id), config.Sub.Fleet, config.Sub.Permanent)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	data, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	entries := strings.Split(bytes.NewBuffer(data).String(), "\n")
	entries = entries[:len(entries)-1] // remove last element

	return entries, nil
}

func Edit(path string, lineNr int) error {
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

func DeleteBuffer() error {
	cmd := exec.Command("/bin/bash", config.Scripts.Clear)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
