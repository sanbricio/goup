package dependency

// Dependency represents a Go module dependency
type Dependency struct {
	Path     string // Module path (e.g., "github.com/gin-gonic/gin")
	Version  string // Current version (e.g., "v1.9.1")
	Indirect bool   // Whether this is an indirect dependency
}

// String returns a string representation of the dependency
func (d Dependency) String() string {
	suffix := ""
	if d.Indirect {
		suffix = " (indirect)"
	}
	return d.Path + "@" + d.Version + suffix
}

// Manager defines the interface for managing Go module dependencies
type Manager interface {
	// GetDependencies reads and parses dependencies from go.mod
	GetDependencies() ([]Dependency, error)
	// FilterDependencies filters dependencies based on criteria
	FilterDependencies(deps []Dependency, includeIndirect bool) []Dependency
}
