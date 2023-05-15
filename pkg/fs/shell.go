package fs

import (
	"os"
	"os/exec"
)

func Exec(command string) error {
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
