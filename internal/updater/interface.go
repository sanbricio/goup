package updater

import "goup/internal/dependency"

// UpdateResult contains the result of an update operation
type UpdateResult struct {
	Updated []dependency.Dependency
	Failed  []UpdateError
	Success bool
}

// UpdateError represents an error that occurred during update
type UpdateError struct {
	Dependency dependency.Dependency
	Error      error
}

// Updater defines the interface for updating dependencies
type Updater interface {
	// UpdateDependencies updates the specified dependencies
	UpdateDependencies(deps []dependency.Dependency, verbose bool) UpdateResult
	// RunModTidy runs go mod tidy to clean up the module
	RunModTidy(verbose bool) error
}

// CommandRunner defines the interface for running system commands
type CommandRunner interface {
	// Run executes a command and returns the result
	Run(name string, args []string, verbose bool) error
}