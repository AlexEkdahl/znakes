# Define the name of the binary
BINARY=game

# Define the names of the server and client binaries
SERVER_BINARY=server
CLIENT_BINARY=client

# Set the default goal to "build"
.DEFAULT_GOAL := build

# Build the server binary
$(SERVER_BINARY):
	@go build -ldflags="-s -w" -o bin/$(SERVER_BINARY) cmd/server/main.go

# Build the client binary
$(CLIENT_BINARY):
	@go build -ldflags="-s -w" -o bin/$(CLIENT_BINARY) cmd/client/main.go

# Build both the server and client binaries
build: $(SERVER_BINARY) $(CLIENT_BINARY)
	@go build -ldflags="-s -w" -o bin/$(BINARY)

# Run the server
run-server: build
	@./bin/$(SERVER_BINARY)

# Run the client
run-client: build
	@./bin/$(CLIENT_BINARY)

# Run the game (server and client)
run: build
	@./bin/$(BINARY)

# Run tests
test:
	go test ./... -v

# Generate Go code from protobuf file
generate:
	go generate ./...

# Clean up
clean:
	rm -f bin/$(BINARY) bin/$(SERVER_BINARY) bin/$(CLIENT_BINARY)
