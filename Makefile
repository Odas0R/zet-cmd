build:
	rm -f ./zet && go build -o zet ./cmd/main
test:
	clear && gotestsum --format testname ./cmd/main
watch:
	find . -name '*.go' | entr -cp richgo test ./cmd/main

.PHONY: build test
