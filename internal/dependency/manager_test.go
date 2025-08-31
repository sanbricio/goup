package dependency

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
		"github.com/gin-gonic/gin",
		"github.com/stretchr/testify",
		"github.com/davecgh/go-spew",
		"github.com/pmezard/go-difflib",
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
		if slices.Contains(indirectDeps, dep.Path) {
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

func TestGetUpdatableDependenciesCommandFails(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create invalid go.mod to make go list fail
	goModContent := `module testmodule

go 1.21

require (
	invalid-module-path v1.0.0
)
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	_, err = manager.GetUpdatableDependencies()

	// Should return an error when go list fails
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check for updates")
}

func TestGetUpdatableDependenciesWithMalformedAndValidData(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create go.mod with actual dependencies
	goModContent := `module testmodule

go 1.21

require (
	github.com/stretchr/testify v1.8.0
)
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	deps, err := manager.GetUpdatableDependencies()

	assert.Error(t, err)
	assert.Nil(t, deps)
}

func TestGetUpdatableDependenciesFilteringLogic(t *testing.T) {

	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	goModContent := `module testmodule

go 1.21

require (
    github.com/stretchr/testify v1.8.0
    golang.org/x/mod v0.10.0
)
`
	require.NoError(t, os.WriteFile(goModPath, []byte(goModContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "stub.go"), []byte("package main"), 0644))

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(originalDir) })
	require.NoError(t, os.Chdir(tempDir))

	cmd := exec.Command("go", "mod", "tidy")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\n%s", err, string(out))
	}

	manager := NewManagerWithPath(goModPath)

	deps, err := manager.GetUpdatableDependencies()
	require.NoError(t, err)

	for _, dep := range deps {
		assert.True(t, dep.HasUpdate, "All returned dependencies should have updates available")
		assert.NotEmpty(t, dep.Path, "Path should not be empty")
		assert.NotEmpty(t, dep.Version, "Current version should not be empty")
		assert.NotEmpty(t, dep.NewVersion, "New version should not be empty")
	}

	if len(deps) == 0 {
		t.Skip("No dependencies with updates found (possibly no network / proxy issue)")
	}

	split := len(deps)
	for i, d := range deps {
		if d.Indirect {
			split = i
			break
		}
	}

	for i := 0; i < split; i++ {
		assert.False(t, deps[i].Indirect, "Direct dependencies should come first")
	}
	for i := split; i < len(deps); i++ {
		assert.True(t, deps[i].Indirect, "Indirect dependencies should come after all direct ones")
	}

	isSorted := func(start, end int) bool {
		for i := start + 1; i < end; i++ {
			if !(deps[i-1].Path < deps[i].Path) {
				return false
			}
		}
		return true
	}
	assert.True(t, isSorted(0, split), "Direct dependencies should be sorted alphabetically")
	assert.True(t, isSorted(split, len(deps)), "Indirect dependencies should be sorted alphabetically")
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

func TestNewManager(t *testing.T) {
	manager := NewManager()

	// We can't easily test the internal goModPath field since it's private,
	// but we can test that the manager was created and works
	assert.NotNil(t, manager)

	// Try to call a method to ensure it was properly initialized
	filtered := manager.FilterDependencies([]Dependency{}, false)
	assert.Empty(t, filtered)
}

func TestNewManagerWithPath(t *testing.T) {
	customPath := "/custom/path/go.mod"
	manager := NewManagerWithPath(customPath)

	assert.NotNil(t, manager)

	// Test that it tries to use the custom path
	_, err := manager.GetDependencies()
	// Should error because the path doesn't exist, but confirms it's using custom path
	assert.Error(t, err)
	assert.Contains(t, err.Error(), customPath)
}

func TestGetDependencies_OnlyIndirectDependencies(t *testing.T) {
	tempDir := t.TempDir()
	goModPath := filepath.Join(tempDir, "go.mod")

	// Create a go.mod with only indirect dependencies
	goModContent := `module testmodule

go 1.21

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)
`

	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	manager := NewManagerWithPath(goModPath)

	deps, err := manager.GetDependencies()
	require.NoError(t, err)

	assert.Len(t, deps, 2)

	// All should be indirect
	for _, dep := range deps {
		assert.True(t, dep.Indirect, "Expected %s to be indirect", dep.Path)
	}
}

func TestFilterDependenciesMixedTypes(t *testing.T) {
	deps := []Dependency{
		{Path: "github.com/direct1", Version: "v1.0.0", Indirect: false},
		{Path: "github.com/indirect1", Version: "v1.0.0", Indirect: true},
		{Path: "github.com/direct2", Version: "v1.0.0", Indirect: false},
		{Path: "github.com/indirect2", Version: "v1.0.0", Indirect: true},
		{Path: "github.com/direct3", Version: "v1.0.0", Indirect: false},
	}

	manager := NewManager()

	t.Run("filter direct only from mixed", func(t *testing.T) {
		filtered := manager.FilterDependencies(deps, false)

		assert.Len(t, filtered, 3)
		expectedDirect := []string{"github.com/direct1", "github.com/direct2", "github.com/direct3"}

		for i, dep := range filtered {
			assert.False(t, dep.Indirect)
			assert.Equal(t, expectedDirect[i], dep.Path)
		}
	})

	t.Run("include all from mixed", func(t *testing.T) {
		filtered := manager.FilterDependencies(deps, true)

		assert.Len(t, filtered, 5)
		assert.Equal(t, deps, filtered)
	})
}
