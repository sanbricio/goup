package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"goup/internal/config"
	"goup/internal/dependency"
	"goup/internal/mocks"
	"goup/internal/selector"
	"goup/internal/updater"
)

func TestRunNoDependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	// Setup expectations
	console.EXPECT().Header().Times(1)
	depMgr.EXPECT().GetUpdatableDependencies().Return([]dependency.Dependency{}, nil).Times(1)
	console.EXPECT().Info("All dependencies are up to date! 🎉").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunGetDependenciesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	// Setup expectations
	console.EXPECT().Header().Times(1)
	depMgr.EXPECT().GetUpdatableDependencies().Return(nil, errors.New("failed to read go.mod")).Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check for dependency updates")
}

func TestRunNoDirectDependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "golang.org/x/crypto", Version: "v0.14.0", Indirect: true},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return([]dependency.Dependency{}).Times(1)
	console.EXPECT().Info("All direct dependencies are up to date! 🎉").Times(1)
	console.EXPECT().Info("(%d indirect dependencies have updates available, use --all to include them)", 1).Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunDryRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{DryRun: true}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	// Setup expectations
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Header().Times(1)
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 1 direct dependencies with available updates:").Times(1)
	console.EXPECT().Warning("Dry run mode - no actual updates will be performed").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunSelectiveModeCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Selective: true}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	sel.EXPECT().Select(deps, false).Return(selector.SelectionResult{Cancelled: true}).Times(1)
	console.EXPECT().Info("No dependencies selected for update").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunSelectiveModeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Selective: true}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	sel.EXPECT().Select(deps, false).Return(selector.SelectionResult{Error: errors.New("selection failed")}).Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dependency selection failed")
}

func TestRunInteractiveModeDeclined(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{Interactive: true}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 1 direct dependencies with available updates:").Times(1)
	console.EXPECT().Confirm("Do you want to proceed with the update?").Return(false).Times(1)
	console.EXPECT().Info("Update cancelled").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunSuccessfulUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	updateResult := updater.UpdateResult{
		Updated: deps,
		Success: true,
		Failed:  []updater.UpdateError{},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 1 direct dependencies with available updates:").Times(1)
	console.EXPECT().Info("Updating dependencies...").Times(1)

	// Progress usa formato con argumentos
	console.EXPECT().Progress("Updating %s... (%d/%d)", "github.com/gin-gonic/gin", 1, 1).Times(1)

	// Expect individual update call
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	console.EXPECT().Success("✓ Updated %s", "github.com/gin-gonic/gin").Times(1)

	// Expect final update call
	upd.EXPECT().UpdateDependencies(deps, false).Return(updateResult).Times(1)
	console.EXPECT().PrintUpdateResult(1, 1, false).Times(1)

	// Expect mod tidy
	console.EXPECT().Info("Running go mod tidy...").Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)
	console.EXPECT().Success("✓ go mod tidy completed").Times(1)
	console.EXPECT().Success("Dependency update completed!").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunUpdateWithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
		{Path: "github.com/bad/package", Version: "v1.0.0", Indirect: false},
	}

	updateResult := updater.UpdateResult{
		Updated: []dependency.Dependency{deps[0]},
		Failed: []updater.UpdateError{
			{Dependency: deps[1], Error: errors.New("update failed")},
		},
		Success: false,
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 2 direct dependencies with available updates:").Times(1)
	console.EXPECT().Info("Updating dependencies...").Times(1)

	// Progress calls for each dependency
	console.EXPECT().Progress("Updating %s... (%d/%d)", "github.com/gin-gonic/gin", 1, 2).Times(1)
	console.EXPECT().Progress("Updating %s... (%d/%d)", "github.com/bad/package", 2, 2).Times(1)

	// Individual updates
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	console.EXPECT().Success("✓ Updated %s", "github.com/gin-gonic/gin").Times(1)

	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[1]}, false).Return(updater.UpdateResult{Success: false}).Times(1)

	// Final update result
	upd.EXPECT().UpdateDependencies(deps, false).Return(updateResult).Times(1)
	console.EXPECT().PrintUpdateResult(1, 2, true).Times(1)
	console.EXPECT().Error("Failed to update %s: %v", "github.com/bad/package", errors.New("update failed")).Times(1)

	// Mod tidy
	console.EXPECT().Info("Running go mod tidy...").Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)
	console.EXPECT().Success("✓ go mod tidy completed").Times(1)
	console.EXPECT().Success("Dependency update completed!").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestRunModTidyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "github.com/gin-gonic/gin", Version: "v1.9.1", Indirect: false},
	}

	updateResult := updater.UpdateResult{
		Updated: deps,
		Success: true,
		Failed:  []updater.UpdateError{},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 1 direct dependencies with available updates:").Times(1)
	console.EXPECT().Info("Updating dependencies...").Times(1)
	console.EXPECT().Progress("Updating %s... (%d/%d)", "github.com/gin-gonic/gin", 1, 1).Times(1)

	// Updates
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	console.EXPECT().Success("✓ Updated %s", "github.com/gin-gonic/gin").Times(1)
	upd.EXPECT().UpdateDependencies(deps, false).Return(updateResult).Times(1)
	console.EXPECT().PrintUpdateResult(1, 1, false).Times(1)

	// Mod tidy fails
	console.EXPECT().Info("Running go mod tidy...").Times(1)
	upd.EXPECT().RunModTidy(false).Return(errors.New("mod tidy failed")).Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go mod tidy failed")
}

// New test for testing the updated case with --all flag
func TestRunNoDirectDependenciesWithAllFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{All: true}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	deps := []dependency.Dependency{
		{Path: "golang.org/x/crypto", Version: "v0.14.0", Indirect: true},
	}

	// Setup expectations
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, true).Return(deps).Times(1)
	console.EXPECT().PrintDependencies(deps, "Found 1 all dependencies with available updates:").Times(1)
	console.EXPECT().Info("Updating dependencies...").Times(1)

	// Progress and update
	console.EXPECT().Progress("Updating %s... (%d/%d)", "golang.org/x/crypto", 1, 1).Times(1)
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	console.EXPECT().Success("✓ Updated %s", "golang.org/x/crypto").Times(1)

	// Final update
	updateResult := updater.UpdateResult{Updated: deps, Success: true, Failed: []updater.UpdateError{}}
	upd.EXPECT().UpdateDependencies(deps, false).Return(updateResult).Times(1)
	console.EXPECT().PrintUpdateResult(1, 1, false).Times(1)

	// Mod tidy
	console.EXPECT().Info("Running go mod tidy...").Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)
	console.EXPECT().Success("✓ go mod tidy completed").Times(1)
	console.EXPECT().Success("Dependency update completed!").Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	console := mocks.NewMockConsole(ctrl)
	depMgr := mocks.NewMockManager(ctrl)
	sel := mocks.NewMockSelector(ctrl)
	upd := mocks.NewMockUpdater(ctrl)

	app := New(cfg, console, depMgr, sel, upd)

	require.NotNil(t, app)
	assert.Equal(t, cfg, app.config)
	assert.Equal(t, console, app.console)
	assert.Equal(t, depMgr, app.depMgr)
	assert.Equal(t, sel, app.selector)
	assert.Equal(t, upd, app.updater)
}
