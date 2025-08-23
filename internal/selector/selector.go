package selector

import (
	"fmt"
	"strconv"
	"strings"

	"goup/internal/dependency"
)

// interactiveSelector implements the Selector interface
type interactiveSelector struct {
	ui     UIInterface
	parser Parser
}

// NewInteractiveSelector creates a new interactive selector
func NewInteractiveSelector(ui UIInterface) Selector {
	return &interactiveSelector{
		ui:     ui,
		parser: NewSelectionParser(),
	}
}

// Select allows the user to choose which dependencies to update
func (s *interactiveSelector) Select(deps []dependency.Dependency, includeIndirect bool) SelectionResult {
	if len(deps) == 0 {
		return SelectionResult{Selected: []dependency.Dependency{}}
	}

	typeStr := "direct"
	if includeIndirect {
		typeStr = "all"
	}

	s.ui.Info("Found %d %s dependencies with available updates:", len(deps), typeStr)
	s.ui.PrintDependencies(deps, "")

	s.showSelectionHelp()

	for {
		input, err := s.ui.ReadInput("Select dependencies to update")
		if err != nil {
			return SelectionResult{Error: fmt.Errorf("reading input: %w", err)}
		}

		input = strings.TrimSpace(input)
		if input == "" {
			return SelectionResult{Cancelled: true}
		}

		selected, err := s.parser.ParseSelection(input, deps)
		if err != nil {
			s.ui.Error("Invalid selection: %v", err)
			continue
		}

		if len(selected) == 0 {
			s.ui.Error("No dependencies matched your selection")
			continue
		}

		// Show selected dependencies and confirm
		s.ui.Success("Selected %d dependencies:", len(selected))
		s.ui.PrintDependencies(selected, "")

		if s.ui.Confirm("Proceed with these selected dependencies?") {
			return SelectionResult{Selected: selected}
		}

		s.ui.Info("Let's try again...")
	}
}

func (s *interactiveSelector) showSelectionHelp() {
	s.ui.Info("Selection options:")
	fmt.Println("  • Enter numbers (e.g., 1,3,5 or 1-3 or 1,3-5)")
	fmt.Println("  • Enter 'all' to select all dependencies")
	fmt.Println("  • Enter package names or patterns (e.g., 'github.com/gin*')")
	fmt.Println("  • Press Enter without input to cancel")
	fmt.Println()
}

// selectionParser implements the Parser interface
type selectionParser struct{}

// NewSelectionParser creates a new selection parser
func NewSelectionParser() Parser {
	return &selectionParser{}
}

// ParseSelection parses user input and returns selected dependencies
func (p *selectionParser) ParseSelection(input string, deps []dependency.Dependency) ([]dependency.Dependency, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "all" {
		return deps, nil
	}

	var selected []dependency.Dependency
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Check if it's a range (e.g., "1-3")
		if strings.Contains(part, "-") {
			rangeDeps, err := p.parseRange(part, deps)
			if err != nil {
				return nil, err
			}
			selected = append(selected, rangeDeps...)
			continue
		}

		// Check if it's a number
		if num, err := strconv.Atoi(part); err == nil {
			if num < 1 || num > len(deps) {
				return nil, fmt.Errorf("number %d is out of range (1-%d)", num, len(deps))
			}
			dep := deps[num-1]
			if !containsDependency(selected, dep) {
				selected = append(selected, dep)
			}
			continue
		}

		// Check if it's a package name or pattern
		matched := p.matchPattern(part, deps, &selected)
		if !matched {
			return nil, fmt.Errorf("no dependencies match pattern: %s", part)
		}
	}

	return selected, nil
}

func (p *selectionParser) parseRange(part string, deps []dependency.Dependency) ([]dependency.Dependency, error) {
	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 {
		return nil, fmt.Errorf("invalid range format: %s", part)
	}

	start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid start number in range: %s", rangeParts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid end number in range: %s", rangeParts[1])
	}

	if start < 1 || end > len(deps) || start > end {
		return nil, fmt.Errorf("range %d-%d is out of bounds (1-%d)", start, end, len(deps))
	}

	var rangeDeps []dependency.Dependency
	for i := start - 1; i < end; i++ {
		rangeDeps = append(rangeDeps, deps[i])
	}

	return rangeDeps, nil
}

func (p *selectionParser) matchPattern(pattern string, deps []dependency.Dependency, selected *[]dependency.Dependency) bool {
	matched := false
	for _, dep := range deps {
		if matchesPattern(strings.ToLower(dep.Path), pattern) {
			if !containsDependency(*selected, dep) {
				*selected = append(*selected, dep)
				matched = true
			}
		}
	}
	return matched
}

// matchesPattern checks if a path matches a pattern with * wildcards
func matchesPattern(path, pattern string) bool {
	// Simple pattern matching with * wildcard
	if !strings.Contains(pattern, "*") {
		return strings.Contains(path, pattern)
	}

	// Split pattern by * and check if all parts are present in order
	parts := strings.Split(pattern, "*")
	index := 0

	for _, part := range parts {
		if part == "" {
			continue
		}
		newIndex := strings.Index(path[index:], part)
		if newIndex == -1 {
			return false
		}
		index += newIndex + len(part)
	}

	return true
}

// containsDependency checks if a dependency is already in the slice
func containsDependency(slice []dependency.Dependency, dep dependency.Dependency) bool {
	for _, item := range slice {
		if item.Path == dep.Path {
			return true
		}
	}
	return false
}
