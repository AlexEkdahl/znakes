build:
	@go build -ldflags="-s -w" -o bin/game

run: build
	@./bin/game

test: build
	go test ./... -v
