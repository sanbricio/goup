package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeToDirectory(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	t.Run("change to valid directory with go.mod", func(t *testing.T) {
		// Create temporary directory with go.mod
		tempDir := t.TempDir()
		goModPath := filepath.Join(tempDir, "go.mod")

		err := os.WriteFile(goModPath, []byte("module test\ngo 1.21\n"), 0644)
		require.NoError(t, err)

		err = changeToDirectory(tempDir)
		assert.NoError(t, err)

		// Verify we're in the correct directory
		currentDir, err := os.Getwd()
		require.NoError(t, err)

		expectedDir, err := filepath.Abs(tempDir)
		require.NoError(t, err)

		assert.Equal(t, expectedDir, currentDir)
	})

	t.Run("error when directory doesn't exist", func(t *testing.T) {
		err := changeToDirectory("/nonexistent/directory")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("error when path is not a directory", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "testfile")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		err = changeToDirectory(tempFile.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not a directory")
	})

	t.Run("error when directory has no go.mod", func(t *testing.T) {
		// Create temporary directory without go.mod
		tempDir := t.TempDir()

		err := changeToDirectory(tempDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no go.mod file found")
		assert.Contains(t, err.Error(), "not a Go module")
	})

	t.Run("handle relative paths", func(t *testing.T) {
		// Create temporary directory with go.mod
		tempDir := t.TempDir()
		goModPath := filepath.Join(tempDir, "go.mod")

		err := os.WriteFile(goModPath, []byte("module test\ngo 1.21\n"), 0644)
		require.NoError(t, err)

		// Change to parent of temp dir first
		parentDir := filepath.Dir(tempDir)
		err = os.Chdir(parentDir)
		require.NoError(t, err)

		// Use relative path
		relativePath := filepath.Base(tempDir)
		err = changeToDirectory(relativePath)
		assert.NoError(t, err)

		// Verify we're in the correct directory
		currentDir, err := os.Getwd()
		require.NoError(t, err)

		expectedDir, err := filepath.Abs(tempDir)
		require.NoError(t, err)

		assert.Equal(t, expectedDir, currentDir)
	})
}

func TestParseFlags(t *testing.T) {
	t.Run("parse directory argument", func(t *testing.T) {
		args := []string{"goup", "/path/to/project"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Equal(t, "/path/to/project", targetDir)
		assert.NotNil(t, config)
	})

	t.Run("parse flags with directory", func(t *testing.T) {
		args := []string{"goup", "--list", "--verbose", "/path/to/project"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Equal(t, "/path/to/project", targetDir)
		assert.True(t, config.List)
		assert.True(t, config.Verbose)
	})

	t.Run("no directory argument", func(t *testing.T) {
		args := []string{"goup", "--interactive"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Empty(t, targetDir)
		assert.True(t, config.Interactive)
	})

	t.Run("parse all flags", func(t *testing.T) {
		args := []string{
			"goup",
			"--list",
			"--interactive",
			"--verbose",
			"--no-color",
			"--all",
			"--select",
		}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Empty(t, targetDir)
		assert.True(t, config.List)
		assert.True(t, config.Interactive)
		assert.True(t, config.Verbose)
		assert.True(t, config.NoColor)
		assert.True(t, config.All)
		assert.True(t, config.Selective)
	})

	t.Run("flags after directory are ignored", func(t *testing.T) {
		// This documents the behavior that flags after non-flag arguments are ignored
		// This is standard Go flag package behavior
		args := []string{"goup", "--list", "/some/path", "--verbose"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Equal(t, "/some/path", targetDir)
		assert.True(t, config.List)
		assert.False(t, config.Verbose) // --verbose is ignored because it comes after the directory
	})

	t.Run("flags before directory only", func(t *testing.T) {
		// Note: Go's flag package stops parsing flags after the first non-flag argument
		// So flags after the directory path won't be parsed
		args := []string{"goup", "--list", "--verbose", "/some/path"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Equal(t, "/some/path", targetDir)
		assert.True(t, config.List)
		assert.True(t, config.Verbose)
	})

	t.Run("only program name", func(t *testing.T) {
		args := []string{"goup"}

		config, targetDir := parseFlagsWithArgs(args)

		assert.Empty(t, targetDir)
		assert.NotNil(t, config)
		assert.False(t, config.List)
		assert.False(t, config.Interactive)
		assert.False(t, config.Verbose)
		assert.False(t, config.NoColor)
		assert.False(t, config.All)
		assert.False(t, config.Selective)
	})
}
