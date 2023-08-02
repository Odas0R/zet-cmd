build:
	TEST=true go test -tags "fts5" ./... \
	&& go build -tags "fts5" -o zet ./cmd/zet \
		&& goose -dir ./migrations sqlite3 ./zettel.db up
watch:
	find . -name '*.go' | entr -cs 'TEST=true go test -tags "fts5" ./... && go build -tags "fts5" -o zet ./cmd/zet'
watch-verbose:
	find . -name '*.go' | entr -cs 'TEST=true go test -v -tags "fts5" ./... && go build -tags "fts5" -o zet ./cmd/zet'
db:
	sqlite3 ./zettel.db
new:
	@read -p "Enter the name of the new migration: " name; \
		goose -dir ./migrations sqlite3 ./zettel.db create $$name sql
up:
	goose -dir ./migrations sqlite3 ./zettel.db up-by-one
down:
	goose -dir ./migrations sqlite3 ./zettel.db down
status:
	goose -dir ./migrations sqlite3 ./zettel.db status
schema:
	sqlite-utils schema zettel.db
test:
	go test -tags "fts5" ./...

.PHONY: build test new up down status schema watch
