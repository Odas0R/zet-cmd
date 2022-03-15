build:
	go build -o zet && ./zet

new:
	@go build -o zet && ./zet new THis is a test title

.PHONY: build
