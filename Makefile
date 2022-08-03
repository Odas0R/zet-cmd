build:
	rm -f ./zet && go build -o zet ./cmd/main
test:
	gotestsum --format testname ./cmd/main
watch:
	gotestsum --watch --hide-summary=skipped --format standard-quiet ./cmd/main

.PHONY: build test
