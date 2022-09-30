package main

import (
	"log"
	"os"
	"testing"
)

const TEST_DIR string = "/tmp/foo"

func TestMain(m *testing.M) {
	// cleanup
	if err := os.RemoveAll("/tmp/foo"); err != nil {
		log.Fatalf("error: failed to cleanup %V", err)
	}

	// initialize config
	if err := config.Init(TEST_DIR); err != nil {
		log.Fatalf("error: failed to initialize config %V", err)
	}
	// initialize history
	if err := history.Init(TEST_DIR, ".history"); err != nil {
		log.Fatalf("error: failed to initialize history %V", err)
	}

	exitVal := m.Run()

	// cleanup
	if err := os.RemoveAll("/tmp/foo"); err != nil {
		log.Fatalf("error: failed to cleanup %V", err)
	}

	os.Exit(exitVal)
}
