package main

import (
	"os/exec"
	"path/filepath"
)

func Query() error {
	query, err := filepath.Abs("./scripts/query")
	if err != nil {
		return err
	}

	// Execute the script query
	cmd := exec.Command("/bin/bash", query)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
