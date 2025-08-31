package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"goup/internal/app"
	"goup/internal/config"
	"goup/internal/dependency"
	"goup/internal/selector"
	"goup/internal/ui"
	"goup/internal/updater"
)

func main() {
	cfg, targetDir := parseFlags()

	// Change to target directory if specified
	if targetDir != "" {
		if err := changeToDirectory(targetDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize dependencies using dependency injection
	console := ui.NewConsole(cfg)
	depManager := dependency.NewManager()
	depSelector := selector.NewInteractiveSelector(console)
	depUpdater := updater.NewGoUpdater()

	// Create and run the application
	application := app.New(cfg, console, depManager, depSelector, depUpdater)

	if err := application.Run(); err != nil {
		console.Error("Application failed: %v", err)
		os.Exit(1)
	}
}

func parseFlags() (*config.Config, string) {
	return parseFlagsWithArgs(os.Args)
}

func parseFlagsWithArgs(args []string) (*config.Config, string) {
	cfg := &config.Config{}

	// Create a new FlagSet to avoid global state issues in tests
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	fs.BoolVar(&cfg.List, "list", false, "List all upgradeable dependencies")
	fs.BoolVar(&cfg.Interactive, "interactive", false, "Ask for confirmation before updating")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Show detailed output")
	fs.BoolVar(&cfg.NoColor, "no-color", false, "Disable colored output")
	fs.BoolVar(&cfg.All, "all", false, "Update indirect dependencies as well")
	fs.BoolVar(&cfg.Selective, "select", false, "Interactively select which dependencies to update")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [directory]\n\n", args[0])
		fmt.Fprintf(os.Stderr, "goup - Go dependency updater\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  directory    Path to Go project directory (default: current directory)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                       		# Update direct dependencies in current directory\n", args[0])
		fmt.Fprintf(os.Stderr, "  %s /path/to/project      		# Update direct dependencies in specified directory\n", args[0])
		fmt.Fprintf(os.Stderr, "  %s /path/to/project --all     # Update direct dependencies in specified directory\n", args[0])
		fmt.Fprintf(os.Stderr, "  %s --select              		# Interactively select dependencies to update\n", args[0])
	}

	// Parse the arguments (skip the program name)
	err := fs.Parse(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Get target directory from command line arguments
	var targetDir string
	if fs.NArg() > 0 {
		targetDir = fs.Arg(0)
	}

	return cfg, targetDir
}

func changeToDirectory(targetDir string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path '%s': %w", targetDir, err)
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", absPath)
	} else if err != nil {
		return fmt.Errorf("failed to access directory '%s': %w", absPath, err)
	}

	// Check if it's actually a directory
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to check directory '%s': %w", absPath, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("'%s' is not a directory", absPath)
	}

	// Change to the directory
	if err := os.Chdir(absPath); err != nil {
		return fmt.Errorf("failed to change to directory '%s': %w", absPath, err)
	}

	// Verify go.mod exists in the target directory
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("no go.mod file found in directory '%s' - not a Go module", absPath)
	}

	fmt.Printf("Working in directory: %s\n", absPath)
	return nil
}
