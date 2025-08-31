package updater

import (
	"fmt"
	"os"
	"os/exec"

	"goup/internal/dependency"
)

// goUpdater implements the Updater interface using Go commands
type goUpdater struct {
	commandRunner CommandRunner
}

// NewGoUpdater creates a new Go updater
func NewGoUpdater() Updater {
	return &goUpdater{
		commandRunner: &systemCommandRunner{},
	}
}

// NewGoUpdaterWithRunner creates a new Go updater with a custom command runner
func NewGoUpdaterWithRunner(runner CommandRunner) Updater {
	return &goUpdater{
		commandRunner: runner,
	}
}

// UpdateDependencies updates the specified dependencies individually
func (u *goUpdater) UpdateDependencies(deps []dependency.Dependency, verbose bool) UpdateResult {
	result := UpdateResult{
		Updated: make([]dependency.Dependency, 0),
		Failed:  make([]UpdateError, 0),
	}

	for _, dep := range deps {
		// Try to update each dependency individually
		// If one fails, add to Failed slice and continue with others
		err := u.commandRunner.Run("go", []string{"get", "-u", dep.Path}, verbose)
		if err != nil {
			result.Failed = append(result.Failed, UpdateError{
				Dependency: dep,
				Error:      err, // Keep original error for better reporting
			})
		} else {
			result.Updated = append(result.Updated, dep)
		}
	}

	result.Success = len(result.Failed) == 0
	return result
}

// RunModTidy runs go mod tidy to clean up the module
func (u *goUpdater) RunModTidy(verbose bool) error {
	return u.commandRunner.Run("go", []string{"mod", "tidy"}, verbose)
}

// systemCommandRunner implements CommandRunner using os/exec
type systemCommandRunner struct{}

// Run executes a command and returns the result
func (r *systemCommandRunner) Run(name string, args []string, verbose bool) error {
	cmd := exec.Command(name, args...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Capture output for non-verbose mode to show only on error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
