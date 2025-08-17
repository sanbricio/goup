package main

import (
	"flag"
	"fmt"
	"os"

	"goup/internal/app"
	"goup/internal/config"
	"goup/internal/dependency"
	"goup/internal/selector"
	"goup/internal/ui"
	"goup/internal/updater"
)

func main() {
	cfg := parseFlags()

	// Initialize dependencies using dependency injection
	console := ui.NewConsole(cfg.NoColor)
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

func parseFlags() *config.Config {
	cfg := &config.Config{}

	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Show what would be updated without making changes")
	flag.BoolVar(&cfg.Interactive, "interactive", false, "Ask for confirmation before updating")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Show detailed output")
	flag.BoolVar(&cfg.NoColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&cfg.All, "all", false, "Update indirect dependencies as well")
	flag.BoolVar(&cfg.Selective, "select", false, "Interactively select which dependencies to update")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "goup - Go dependency updater\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --select              # Interactively select dependencies\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --select --all        # Select from all dependencies (including indirect)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --select --dry-run    # Preview selective updates\n", os.Args[0])
	}

	flag.Parse()
	return cfg
}
