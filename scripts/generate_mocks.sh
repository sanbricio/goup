#!/bin/bash

# Script to generate all mocks
echo "Generating mocks..."

# Create mocks directory if it doesn't exist
mkdir -p internal/mocks

# Find mockgen executable (it should be available now thanks to the Makefile)
MOCKGEN=""
if command -v mockgen &> /dev/null; then
    MOCKGEN="mockgen"
else
    echo "❌ mockgen still not found in PATH"
    exit 1
fi

echo "Using mockgen: $(which $MOCKGEN)"

# Generate mock for dependency.Manager
echo "Generating mock for dependency.Manager..."
$MOCKGEN -source=internal/dependency/interface.go -destination=internal/mocks/mock_dependency_manager.go -package=mocks

# Generate mock for ui.Console
echo "Generating mock for ui.Console..."
$MOCKGEN -source=internal/ui/interface.go -destination=internal/mocks/mock_console.go -package=mocks

# Generate mock for selector.Selector
echo "Generating mock for selector.Selector..."
$MOCKGEN -source=internal/selector/interface.go -destination=internal/mocks/mock_selector.go -package=mocks

# Generate mock for updater.Updater
echo "Generating mock for updater.Updater..."
$MOCKGEN -source=internal/updater/interface.go -destination=internal/mocks/mock_updater.go -package=mocks

echo "✅ Mocks generated successfully in internal/mocks/"