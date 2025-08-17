package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	// ANSI color codes
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

type Config struct {
	DryRun      bool
	Interactive bool
	Verbose     bool
	NoColor     bool
	All         bool // Update indirect dependencies too
}

type Dependency struct {
	Path     string
	Version  string
	Indirect bool
}

func main() {
	config := parseFlags()

	printHeader(config)

	deps, err := getDependencies()
	if err != nil {
		printError("Failed to read dependencies: %v", err)
		os.Exit(1)
	}

	if len(deps) == 0 {
		printInfo("No dependencies found to update")
		return
	}

	directDeps := filterDependencies(deps, config.All)
	if len(directDeps) == 0 {
		printInfo("No direct dependencies found to update")
		return
	}

	printDependencies(directDeps, config)

	if config.Interactive && !confirmUpdate() {
		printInfo("Update cancelled")
		return
	}

	if config.DryRun {
		printWarning("Dry run mode - no actual updates will be performed")
		return
	}

	updateDependencies(directDeps, config)
	runModTidy(config)

	printSuccess("Dependency update completed!")
}

func parseFlags() Config {
	var config Config

	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be updated without making changes")
	flag.BoolVar(&config.Interactive, "interactive", false, "Ask for confirmation before updating")
	flag.BoolVar(&config.Verbose, "verbose", false, "Show detailed output")
	flag.BoolVar(&config.NoColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&config.All, "all", false, "Update indirect dependencies as well")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "goup - Go dependency updater\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	return config
}

func printHeader(config Config) {
	if config.NoColor {
		fmt.Println("=== goup - Go Dependency Updater ===")
	} else {
		fmt.Printf("%s%s=== goup - Go Dependency Updater ===%s\n", ColorBold, ColorCyan, ColorReset)
	}
	fmt.Println()
}

func getDependencies() ([]Dependency, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return nil, fmt.Errorf("reading go.mod: %w", err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing go.mod: %w", err)
	}

	var deps []Dependency
	for _, req := range f.Require {
		deps = append(deps, Dependency{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}

	// Sort dependencies alphabetically
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Path < deps[j].Path
	})

	return deps, nil
}

func filterDependencies(deps []Dependency, includeIndirect bool) []Dependency {
	var filtered []Dependency
	for _, dep := range deps {
		if !dep.Indirect || includeIndirect {
			filtered = append(filtered, dep)
		}
	}
	return filtered
}

func printDependencies(deps []Dependency, config Config) {
	typeStr := "direct"
	if config.All {
		typeStr = "all"
	}

	printInfo("Found %d %s dependencies to update:", len(deps), typeStr)
	fmt.Println()

	for i, dep := range deps {
		prefix := fmt.Sprintf("  %d.", i+1)
		status := ""
		if dep.Indirect {
			status = " (indirect)"
		}

		if config.NoColor {
			fmt.Printf("%s %s@%s%s\n", prefix, dep.Path, dep.Version, status)
		} else {
			color := ColorGreen
			if dep.Indirect {
				color = ColorYellow
			}
			fmt.Printf("  %s%d.%s %s%s@%s%s%s%s\n",
				ColorBold, i+1, ColorReset,
				color, dep.Path, dep.Version, ColorReset,
				ColorPurple, status)
		}
	}
	fmt.Println()
}

func confirmUpdate() bool {
	fmt.Print("Do you want to proceed with the update? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func updateDependencies(deps []Dependency, config Config) {
	printInfo("Updating dependencies...")
	fmt.Println()

	errorCount := 0

	for i, dep := range deps {
		printProgress("Updating %s... (%d/%d)", dep.Path, i+1, len(deps))

		cmd := exec.Command("go", "get", "-u", dep.Path)

		if config.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			// Capture output to show only on error
			output, err := cmd.CombinedOutput()
			if err != nil {
				printError("Failed to update %s: %v", dep.Path, err)
				if len(output) > 0 {
					fmt.Printf("Output: %s\n", string(output))
				}
				errorCount++
				continue
			}
		}

		if err := cmd.Run(); err != nil && !config.Verbose {
			printError("Failed to update %s: %v", dep.Path, err)
			errorCount++
		} else {
			printSuccess("✓ Updated %s", dep.Path)
		}
	}

	fmt.Println()
	if errorCount > 0 {
		printWarning("Completed with %d errors", errorCount)
	} else {
		printSuccess("All dependencies updated successfully!")
	}
}

func runModTidy(config Config) {
	printInfo("Running go mod tidy...")

	cmd := exec.Command("go", "mod", "tidy")

	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		printError("Failed to run go mod tidy: %v", err)
	} else {
		printSuccess("✓ go mod tidy completed")
	}
	fmt.Println()
}

// Utility functions for colored output
func printInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[INFO] %s\n", message)
}

func printSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s[SUCCESS]%s %s\n", ColorGreen, ColorReset, message)
}

func printWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s[WARNING]%s %s\n", ColorYellow, ColorReset, message)
}

func printError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, message)
}

func printProgress(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s[PROGRESS]%s %s\n", ColorBlue, ColorReset, message)
}
