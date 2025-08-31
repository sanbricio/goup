package dependency

// Dependency represents a Go module dependency with update information
type Dependency struct {
	Path       string // Module path (e.g., "github.com/gin-gonic/gin")
	Version    string // Current version (e.g., "v1.9.1")
	NewVersion string // Available new version (e.g., "v1.9.2")
	Indirect   bool   // Whether this is an indirect dependency
	HasUpdate  bool   // Whether an update is available
}

// String returns a string representation of the dependency
func (d Dependency) String() string {
	suffix := ""
	if d.Indirect {
		suffix = " (indirect)"
	}

	if d.HasUpdate && d.NewVersion != "" {
		return d.Path + "@" + d.Version + " → " + d.NewVersion + suffix
	}
	return d.Path + "@" + d.Version + suffix
}

func (d Dependency) VersionInfo() string {
	if d.HasUpdate && d.NewVersion != "" {
		return d.Version + " → " + d.NewVersion
	}
	return d.Version
}

// Manager defines the interface for managing Go module dependencies
type Manager interface {
	// GetDependencies reads and parses dependencies from go.mod
	GetDependencies() ([]Dependency, error)
	// FilterDependencies filters dependencies based on criteria
	FilterDependencies(deps []Dependency, includeIndirect bool) []Dependency
	// GetUpdatableDependencies returns only dependencies that have updates available
	GetUpdatableDependencies() ([]Dependency, error)
}
