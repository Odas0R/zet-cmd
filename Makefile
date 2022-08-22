build:
	rm -f ./zet && go build -o zet ./cmd/main
test:
	clear && gotestsum --format standard-verbose ./cmd/main
watch:
	find . -name '*.go' | entr -cp richgo test ./cmd/main

.PHONY: build test
