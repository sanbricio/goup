package dependency

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDependencies(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create a sample go.mod file
	goModContent := `module testmodule

go 1.21

require (
	github.com/stretchr/testify v1.8.0
	github.com/gin-gonic/gin v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	deps, err := manager.GetDependencies()
	require.NoError(t, err)

	// Should have 5 total dependencies (2 direct + 3 indirect)
	assert.Len(t, deps, 5)

	// Check that dependencies are sorted alphabetically
	expectedOrder := []string{
		"github.com/davecgh/go-spew",
		"github.com/gin-gonic/gin",
		"github.com/pmezard/go-difflib",
		"github.com/stretchr/testify",
		"gopkg.in/yaml.v3",
	}

	for i, dep := range deps {
		assert.Equal(t, expectedOrder[i], dep.Path)
	}

	// Check indirect flags
	indirectDeps := []string{
		"github.com/davecgh/go-spew",
		"github.com/pmezard/go-difflib",
		"gopkg.in/yaml.v3",
	}

	for _, dep := range deps {
		if contains(indirectDeps, dep.Path) {
			assert.True(t, dep.Indirect, "Expected %s to be indirect", dep.Path)
		} else {
			assert.False(t, dep.Indirect, "Expected %s to be direct", dep.Path)
		}
	}
}

func TestFilterDependencies(t *testing.T) {
	deps := []Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
		{Path: "golang.org/x/crypto", Version: "v0.14.0", Indirect: true},
		{Path: "github.com/stretchr/testify", Version: "v1.8.4", Indirect: false},
		{Path: "gopkg.in/yaml.v3", Version: "v3.0.1", Indirect: true},
	}

	manager := NewManager()

	t.Run("include only direct dependencies", func(t *testing.T) {
		filtered := manager.FilterDependencies(deps, false)

		assert.Len(t, filtered, 2)
		for _, dep := range filtered {
			assert.False(t, dep.Indirect)
		}
	})

	t.Run("include all dependencies", func(t *testing.T) {
		filtered := manager.FilterDependencies(deps, true)

		assert.Len(t, filtered, 4)
		assert.Equal(t, deps, filtered)
	})

	t.Run("empty dependencies", func(t *testing.T) {
		filtered := manager.FilterDependencies([]Dependency{}, false)
		assert.Empty(t, filtered)
	})
}

func TestGetUpdatableDependenciesInIsolatedEnv(t *testing.T) {
	// Create a completely isolated test environment
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create a minimal go.mod file
	goModContent := `module testmodule

go 1.21
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// Change to the temp directory so go commands work in isolation
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	updatableDeps, err := manager.GetUpdatableDependencies()

	assert.NoError(t, err)
	assert.Empty(t, updatableDeps, "Should have no updatable dependencies in empty module")
}

func TestGetDependenciesInvalidGoMod(t *testing.T) {
	manager := NewManagerWithPath("nonexistent.mod")

	deps, err := manager.GetDependencies()

	assert.Error(t, err)
	assert.Nil(t, deps)
	assert.Contains(t, err.Error(), "reading nonexistent.mod")
}

func TestGetDependenciesMalformedGoMod(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create a malformed go.mod file
	goModContent := `this is not valid go.mod syntax`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	deps, err := manager.GetDependencies()

	assert.Error(t, err)
	assert.Nil(t, deps)
	assert.Contains(t, err.Error(), "parsing")
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
