#!/bin/bash

go mod download

cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

