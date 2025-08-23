package dependency

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// manager implements the Manager interface
type manager struct {
	goModPath string
}

// NewManager creates a new dependency manager
func NewManager() Manager {
	return &manager{
		goModPath: "go.mod",
	}
}

// NewManagerWithPath creates a new dependency manager with a custom go.mod path
func NewManagerWithPath(path string) Manager {
	return &manager{
		goModPath: path,
	}
}

// GetDependencies reads and parses dependencies from go.mod
func (m *manager) GetDependencies() ([]Dependency, error) {
	data, err := os.ReadFile(m.goModPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", m.goModPath, err)
	}

	f, err := modfile.Parse(m.goModPath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", m.goModPath, err)
	}

	var deps []Dependency
	for _, req := range f.Require {
		deps = append(deps, Dependency{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}

	// Sort dependencies alphabetically for consistent output
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Path < deps[j].Path
	})

	return deps, nil
}

// FilterDependencies filters dependencies based on criteria
func (m *manager) FilterDependencies(deps []Dependency, includeIndirect bool) []Dependency {
	if includeIndirect {
		return deps
	}

	var filtered []Dependency
	for _, dep := range deps {
		if !dep.Indirect {
			filtered = append(filtered, dep)
		}
	}
	return filtered
}

// GetUpdatableDependencies returns only dependencies that have updates available
func (m *manager) GetUpdatableDependencies() ([]Dependency, error) {
	// Use 'go list -u -m all' to get ALL dependencies (direct and indirect) with their update info
	cmd := exec.Command("go", "list", "-u", "-m", "-json", "all")
	output, err := cmd.Output()
	if err != nil {
		// If the command fails, fall back to reading from go.mod
		return m.GetDependencies()
	}

	var allDeps []Dependency
	updatableMap := make(map[string]bool)

	decoder := json.NewDecoder(strings.NewReader(string(output)))

	for decoder.More() {
		var module struct {
			Path     string `json:"Path"`
			Version  string `json:"Version"`
			Indirect bool   `json:"Indirect"`
			Main     bool   `json:"Main"`
			Update   *struct {
				Path    string `json:"Path"`
				Version string `json:"Version"`
			} `json:"Update"`
		}

		if err := decoder.Decode(&module); err != nil {
			continue
		}

		// Skip the main module
		if module.Main {
			continue
		}

		// Skip modules without a version (usually means they're local or replaced)
		if module.Version == "" {
			continue
		}

		// Create dependency object
		dep := Dependency{
			Path:     module.Path,
			Version:  module.Version,
			Indirect: module.Indirect,
		}

		allDeps = append(allDeps, dep)

		// If Update field is present, this dependency has an update available
		if module.Update != nil {
			updatableMap[module.Path] = true
		}
	}

	// Filter to only include dependencies with updates
	var updatable []Dependency
	for _, dep := range allDeps {
		if updatableMap[dep.Path] {
			updatable = append(updatable, dep)
		}
	}

	// Sort dependencies alphabetically for consistent output
	sort.Slice(updatable, func(i, j int) bool {
		return updatable[i].Path < updatable[j].Path
	})

	return updatable, nil
}
