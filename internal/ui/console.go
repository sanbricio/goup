package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"goup/internal/dependency"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// console implements the Console interface
type console struct {
	noColor bool
	reader  *bufio.Reader
}

// NewConsole creates a new console UI
func NewConsole(noColor bool) Console {
	return &console{
		noColor: noColor,
		reader:  bufio.NewReader(os.Stdin),
	}
}

// Header displays the application header
func (c *console) Header() {
	if c.noColor {
		fmt.Println("=== goup - Go Dependency Updater ===")
	} else {
		fmt.Printf("%s%s=== goup - Go Dependency Updater ===%s\n", ColorBold, ColorCyan, ColorReset)
	}
	fmt.Println()
}

// Info displays an informational message
func (c *console) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[INFO] %s\n", message)
}

// Success displays a success message
func (c *console) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if c.noColor {
		fmt.Printf("[SUCCESS] %s\n", message)
	} else {
		fmt.Printf("%s[SUCCESS]%s %s\n", ColorGreen, ColorReset, message)
	}
}

// Warning displays a warning message
func (c *console) Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if c.noColor {
		fmt.Printf("[WARNING] %s\n", message)
	} else {
		fmt.Printf("%s[WARNING]%s %s\n", ColorYellow, ColorReset, message)
	}
}

// Error displays an error message
func (c *console) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if c.noColor {
		fmt.Printf("[ERROR] %s\n", message)
	} else {
		fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, message)
	}
}

// Progress displays a progress message
func (c *console) Progress(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if c.noColor {
		fmt.Printf("[PROGRESS] %s\n", message)
	} else {
		fmt.Printf("%s[PROGRESS]%s %s\n", ColorBlue, ColorReset, message)
	}
}

// ReadInput reads a line of input from the user with a prompt
func (c *console) ReadInput(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

// Confirm asks the user for yes/no confirmation
func (c *console) Confirm(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// PrintDependencies displays a numbered list of dependencies
func (c *console) PrintDependencies(deps []dependency.Dependency, title string) {
	if title != "" {
		c.Info("%s", title)
	}
	fmt.Println()

	for i, dep := range deps {
		c.printDependency(i+1, dep)
	}
	fmt.Println()
}

func (c *console) printDependency(index int, dep dependency.Dependency) {
	prefix := fmt.Sprintf("  %d.", index)
	status := ""
	if dep.Indirect {
		status = " (indirect)"
	}

	if c.noColor {
		fmt.Printf("%s %s@%s%s\n", prefix, dep.Path, dep.Version, status)
	} else {
		color := ColorGreen
		if dep.Indirect {
			color = ColorYellow
		}
		fmt.Printf("  %s%d.%s %s%s@%s%s%s%s\n",
			ColorBold, index, ColorReset,
			color, dep.Path, dep.Version, ColorReset,
			ColorPurple, status)
	}
}

// PrintUpdateResult displays the result of an update operation
func (c *console) PrintUpdateResult(updated, total int, hasErrors bool) {
	if hasErrors {
		c.Warning("Completed with %d/%d dependencies updated", updated, total)
	} else {
		c.Success("All %d dependencies updated successfully!", total)
	}
}
