build:
	go build -o zet ./cmd/zet && ./zet
db:
	sqlite3 ./zettel.db
new:
	@read -p "Enter the name of the new zettel: " name; \
		goose -dir ./migrations sqlite3 ./zettel.db create $$name sql
up:
	goose -dir ./migrations sqlite3 ./zettel.db up 1
down:
	goose -dir ./migrations sqlite3 ./zettel.db down 1
status:
	goose -dir ./migrations sqlite3 ./zettel.db status
schema:
	sqlite-utils schema zettel.db
test:
	go test ./...

.PHONY: build test new up down status schema
