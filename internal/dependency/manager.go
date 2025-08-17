package dependency

import (
	"fmt"
	"os"
	"sort"

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
