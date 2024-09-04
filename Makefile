GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.22.0
DB_PATH := ./zettel.db
TMP_DB_PATH := /tmp/zet/zettel.db
MIGRATIONS_DIR := ./migrations

.PHONY: build build-tmp watch new up down status schema test

build:
	go build -tags "fts5" -o zet ./cmd/zet
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up

install:
		sudo install -m 755 ./zet ~/.local/bin/zet

build-tmp:
	TEST=true go test -tags "fts5" ./...
	go build -tags "fts5" -o zet ./cmd/zet
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(TMP_DB_PATH) up

watch:
	find . -name '*.go' | entr -cs 'TEST=true go test -tags "fts5" ./... && go build -tags "fts5" -o zet ./cmd/zet'

new:
	@read -p "Enter the name of the new migration: " name; \
		$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) create $$name sql

up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up-by-one

down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) down

redo:
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) redo

status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) status

schema:
	sqlite3 $(DB_PATH) .schema

test:
	TEST=true go test -tags "fts5" ./...
