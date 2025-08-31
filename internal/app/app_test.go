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
	assert.Contains(t, err.Error(), "failed to read go.mod")
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

func TestRunList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{List: true}
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

	// Setup expectations - Solo UI con AnyTimes
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Info(gomock.Any()).AnyTimes()
	console.EXPECT().Progress(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().ProgressBar(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Success(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintDependencies(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintUpdateResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)

	// Solo la llamada individual (eliminamos la final)
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)

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

	// Setup expectations - Solo UI con AnyTimes
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Info(gomock.Any()).AnyTimes()
	console.EXPECT().Progress(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().ProgressBar(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Success(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintDependencies(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintUpdateResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)

	// Solo llamadas individuales (eliminamos la final)
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[1]}, false).Return(updater.UpdateResult{Success: false}).Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)

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

	// Setup expectations - UI calls with flexibility
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().ProgressBar(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Success(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintDependencies(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintUpdateResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	console.EXPECT().Warning("go mod tidy failed: %v", errors.New("mod tidy failed")).Times(1)

	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, false).Return(deps).Times(1)

	// Individual dependency update succeeds
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{
		Updated: deps,
		Failed:  []updater.UpdateError{},
		Success: true,
	}).Times(1)

	// Mod tidy fails
	upd.EXPECT().RunModTidy(false).Return(errors.New("mod tidy failed")).Times(1)

	app := New(cfg, console, depMgr, sel, upd)
	err := app.Run()

	// ⭐ NEW: Should NOT error - resilient behavior continues even if mod tidy fails
	assert.NoError(t, err)
}

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

	// Setup expectations - Solo UI con AnyTimes
	console.EXPECT().Header().Times(1)
	console.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Info(gomock.Any()).AnyTimes()
	console.EXPECT().Progress(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().ProgressBar(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().Success(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintDependencies(gomock.Any(), gomock.Any()).AnyTimes()
	console.EXPECT().PrintUpdateResult(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	depMgr.EXPECT().GetUpdatableDependencies().Return(deps, nil).Times(1)
	depMgr.EXPECT().FilterDependencies(deps, true).Return(deps).Times(1)

	// Solo llamada individual (eliminamos la final)
	upd.EXPECT().UpdateDependencies([]dependency.Dependency{deps[0]}, false).Return(updater.UpdateResult{Success: true}).Times(1)
	upd.EXPECT().RunModTidy(false).Return(nil).Times(1)

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
