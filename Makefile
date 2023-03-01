setup:
	@echo "Setting up the environment"
	@./scripts/setup.sh

cibuild:
	./scripts/cibuild.sh

#####################################

SERVER_BINARY=znakes
SERVER_SRC=./cmd/server/main.go
BIN_DIR=./bin
.DEFAULT_GOAL := run

$(SERVER_BINARY):
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$@ $(SERVER_SRC) >/dev/null

run: $(SERVER_BINARY)
	$(BIN_DIR)/$(SERVER_BINARY)

test:
	go test ./... -v

clean:
	go clean
	rm -rf $(BIN_DIR)

#####################################

PLATFORMS := linux/amd64 windows/amd64 darwin/amd64
LINUX_AMD64_DIR=./bin/linux-amd64
LINUX_ARM_DIR=./bin/linux-arm
WINDOWS_AMD64_DIR=./bin/windows-amd64
ARM64_DIR=./bin/darwin-arm64

$(SERVER_BINARY)-rasp:
	@GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o $(LINUX_ARM_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null

$(SERVER_BINARY)-linux:
	@GOOS=$(word 1, $(LINUX_AMD64)) GOARCH=$(word 2, $(LINUX_AMD64)) go build -ldflags="-s -w" -o $(LINUX_AMD64_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null

$(SERVER_BINARY)-windows:
	@GOOS=$(word 1, $(WINDOWS_AMD64)) GOARCH=$(word 2, $(WINDOWS_AMD64)) go build -ldflags="-s -w" -o $(WINDOWS_AMD64_DIR)/$(SERVER_BINARY).exe $(SERVER_SRC) >/dev/null

$(SERVER_BINARY)-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(ARM64_DIR)/$(SERVER_BINARY) $(SERVER_SRC) >/dev/null
