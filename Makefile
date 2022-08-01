build:
	rm -f ./zet && go build -o zet ./cmd/main
test:
	go test ./cmd/main -v

.PHONY: build test
