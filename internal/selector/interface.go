package selector

import "goup/internal/dependency"

// SelectionResult contains the result of a dependency selection
type SelectionResult struct {
	Selected  []dependency.Dependency
	Cancelled bool
	Error     error
}

// Selector defines the interface for selecting dependencies
type Selector interface {
	// Select allows the user to choose which dependencies to update
	Select(deps []dependency.Dependency, includeIndirect bool) SelectionResult
}

// Parser defines the interface for parsing selection input
type Parser interface {
	// ParseSelection parses user input and returns selected dependencies
	ParseSelection(input string, deps []dependency.Dependency) ([]dependency.Dependency, error)
}

// UIInterface defines the interface for user interaction
type UIInterface interface {
	// Info displays an informational message
	Info(format string, args ...interface{})

	// Success displays a success message
	Success(format string, args ...interface{})

	// Error displays an error message
	Error(format string, args ...interface{})

	// ReadInput reads a line of input from the user
	ReadInput(prompt string) (string, error)

	// Confirm asks the user for yes/no confirmation
	Confirm(message string) bool

	// PrintDependencies displays a list of dependencies
	PrintDependencies(deps []dependency.Dependency, title string)
}
