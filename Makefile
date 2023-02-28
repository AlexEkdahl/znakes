setup:
	@echo "Setting up the environment"
	@./scripts/setup.sh

cibuild:
	./scripts/cibuild.sh

#####################################

# Define the name of the binary
BINARY=game

# Define the names of the server and client binaries
SERVER_BINARY=server
CLIENT_BINARY=client

# Define the source files and directories
SERVER_SRC=./cmd/server/main.go
CLIENT_SRC=./cmd/client/main.go

# Define the output directories
BIN_DIR=./bin

# Set the default goal to "all"
.DEFAULT_GOAL := all

# Build the server binary silently
$(SERVER_BINARY):
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$@ $(SERVER_SRC) >/dev/null

# Build the client binary silently
$(CLIENT_BINARY):
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$@ $(CLIENT_SRC) >/dev/null

# Build both the server and client binaries silently
all: $(SERVER_BINARY) $(CLIENT_BINARY)
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY) >/dev/null

# Run the server
run-server: $(SERVER_BINARY)
	$(BIN_DIR)/$(SERVER_BINARY)

# Run the client
run-client: $(CLIENT_BINARY)
	$(BIN_DIR)/$(CLIENT_BINARY)

# Run the game (server and client)
run: all
	$(BIN_DIR)/$(BINARY)

# Run tests
test:
	go test ./... -v

# Generate Go code from protobuf file
generate:
	go generate ./...

# Clean up
clean:
	rm -rf $(BIN_DIR)

#####################################

# Define the target platforms
PLATFORMS := linux/amd64 windows/amd64 darwin/amd64
# Define the output directories

# Define the target platform output directories
LINUX_AMD64_DIR=./bin/linux-amd64
LINUX_ARM_DIR=./bin/linux-arm
WINDOWS_AMD64_DIR=./bin/windows-amd64
ARM64_DIR=./bin/darwin-arm64

# Build the server binary for Linux rasp
$(SERVER_BINARY)-rasp:
	@GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o $(LINUX_ARM_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null

# Build the server binary for Linux
$(SERVER_BINARY)-linux:
	@GOOS=$(word 1, $(LINUX_AMD64)) GOARCH=$(word 2, $(LINUX_AMD64)) go build -ldflags="-s -w" -o $(LINUX_AMD64_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null
# Build the server binary for Windows
#
$(SERVER_BINARY)-windows:
	@GOOS=$(word 1, $(WINDOWS_AMD64)) GOARCH=$(word 2, $(WINDOWS_AMD64)) go build -ldflags="-s -w" -o $(WINDOWS_AMD64_DIR)/$(SERVER_BINARY).exe $(SERVER_SRC) >/dev/null

# Build the server binary for Apple Silicon (arm64)
$(SERVER_BINARY)-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(ARM64_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null

# Build the client binary for Linux
$(CLIENT_BINARY)-linux:
	@GOOS=$(word 1, $(LINUX_AMD64)) GOARCH=$(word 2, $(LINUX_AMD64)) go build -ldflags="-s -w" -o $(LINUX_AMD64_DIR)/$(CLIENT_BINARY) $(CLIENT_SRC) >/dev/null

# Build the client binary for Windows
$(CLIENT_BINARY)-windows:
	@GOOS=$(word 1, $(WINDOWS_AMD64)) GOARCH=$(word 2, $(WINDOWS_AMD64)) go build -ldflags="-s -w" -o $(WINDOWS_AMD64_DIR)/$(CLIENT_BINARY).exe $(CLIENT_SRC) >/dev/null

# Build the client binary for Apple Silicon (arm64)
$(CLIENT_BINARY)-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(ARM64_DIR)/$(CLIENT_BINARY) $(CLIENT_SRC) >/dev/null
