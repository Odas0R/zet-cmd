build:
	go build -o zet ./cmd/zet && ./zet
new:
	@read -p "Enter the name of the new zettel: " name; \
		goose -dir ./pkg/database/migrations sqlite3 ./zettel.db create $$name sql
up:
	goose -dir ./pkg/database/migrations sqlite3 ./zettel.db up 1
down:
	goose -dir ./pkg/database/migrations sqlite3 ./zettel.db down 1
status:
	goose -dir ./pkg/database/migrations sqlite3 ./zettel.db status
schema:
	sqlite-utils schema zettel.db
test:
	go test ./...

.PHONY: build test new up down status schema
