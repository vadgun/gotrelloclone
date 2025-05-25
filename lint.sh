#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Find all go.mod files, excluding those in frontend or node_modules
find . -name go.mod -not -path "./frontend/*" -not -path "./node_modules/*" | while read -r modfile; do
    dir=$(dirname "$modfile")
    echo "Running golangci-lint in $dir"
    (cd "$dir" && $HOME/go/bin/golangci-lint run)
done

echo "golangci-lint completed successfully."
