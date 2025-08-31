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
	m.sortDependencies(deps)

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

// GetUpdatableDependencies returns ONLY dependencies that have updates available
func (m *manager) GetUpdatableDependencies() ([]Dependency, error) {
	// Use 'go list -u -m all' to get ALL dependencies with their update info
	out, err := exec.Command("go", "list", "-u", "-m", "-json", "all").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %v\noutput:\n%s", err, string(out))
	}

	var updatableDeps []Dependency
	decoder := json.NewDecoder(strings.NewReader(string(out)))

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

		// Skip modules without a version
		if module.Version == "" {
			continue
		}

		// ONLY add dependencies that have updates available
		if module.Update != nil {
			dep := Dependency{
				Path:       module.Path,
				Version:    module.Version,
				NewVersion: module.Update.Version,
				Indirect:   module.Indirect,
				HasUpdate:  true,
			}
			updatableDeps = append(updatableDeps, dep)
		}
	}

	// Sort dependencies: first direct (alphabetically), then indirect (alphabetically)
	m.sortDependencies(updatableDeps)

	return updatableDeps, nil
}

func (m *manager) sortDependencies(deps []Dependency) {
	sort.Slice(deps, func(i, j int) bool {
		depA, depB := deps[i], deps[j]

		// If one is direct and the other is indirect, direct comes first
		if depA.Indirect != depB.Indirect {
			return !depA.Indirect // !false = true (direct first), !true = false (indirect second)
		}

		// If both are the same type (both direct or both indirect), sort alphabetically
		return depA.Path < depB.Path
	})
}
