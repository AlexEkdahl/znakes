#!/bin/bash

go mod download
install github.com/golangci/golangci-lint/cmd/golangci-lint
go install github.com/segmentio/golines@latest

cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

