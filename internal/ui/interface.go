package ui

import "goup/internal/dependency"

// Console defines the interface for console-based user interaction
type Console interface {
	// Header displays the application header
	Header()

	// Info displays an informational message
	Info(format string, args ...interface{})

	// Success displays a success message
	Success(format string, args ...interface{})

	// Warning displays a warning message
	Warning(format string, args ...interface{})

	// Error displays an error message
	Error(format string, args ...interface{})

	// Progress displays a progress message
	Progress(format string, args ...interface{})

	// ReadInput reads a line of input from the user with a prompt
	ReadInput(prompt string) (string, error)

	// Confirm asks the user for yes/no confirmation
	Confirm(message string) bool

	// PrintDependencies displays a numbered list of dependencies
	PrintDependencies(deps []dependency.Dependency, title string)

	// PrintUpdateResult displays the result of an update operation
	PrintUpdateResult(updated, total int, hasErrors bool)
}
