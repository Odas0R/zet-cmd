build:
	go build -tags "fts5" -o zet ./cmd/zet \
		&& goose -dir ./migrations sqlite3 ./zettel.db up
install:
		sudo install -m 755 ./zet ~/.local/bin/zet
build-tmp:
	TEST=true go test -tags "fts5" ./... \
	&& go build -tags "fts5" -o zet ./cmd/zet \
		&& goose -dir ./migrations sqlite3 /tmp/zet/zettel.db up
watch:
	find . -name '*.go' | entr -cs 'TEST=true go test -tags "fts5" ./... && go build -tags "fts5" -o zet ./cmd/zet'
new:
	@read -p "Enter the name of the new migration: " name; \
		goose -dir ./migrations sqlite3 ./zettel.db create $$name sql
up:
	goose -dir ./migrations sqlite3 ./zettel.db up-by-one
down:
	goose -dir ./migrations sqlite3 ./zettel.db down
redo:
	goose -dir ./migrations sqlite3 ./zettel.db redo
status:
	goose -dir ./migrations sqlite3 ./zettel.db status
schema:
	sqlite3 ./zettel.db .schema
test:
	TEST=true go test -tags "fts5" ./...

.PHONY: build test new up down status schema watch install
