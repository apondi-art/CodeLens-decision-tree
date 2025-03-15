#!/bin/bash

# Set the working directory to the script's location
cd "$(dirname "$0")" || exit

echo "Formatting Go files..."
gofmt -s -w .
command -v goimports >/dev/null 2>&1 && goimports -w .

echo "Removing old binary..."
rm -f bin/dt  # Deletes the old binary if it exists

echo "Building project..."
mkdir -p bin
go build -o bin/dt ./cmd/dt

echo "Build complete. Executable located at bin/dt"
