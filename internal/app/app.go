package app

import (
	"fmt"

	"goup/internal/config"
	"goup/internal/dependency"
	"goup/internal/selector"
	"goup/internal/ui"
	"goup/internal/updater"
)

// App represents the main application
type App struct {
	config   *config.Config
	console  ui.Console
	depMgr   dependency.Manager
	selector selector.Selector
	updater  updater.Updater
}

// New creates a new application instance
func New(
	cfg *config.Config,
	console ui.Console,
	depMgr dependency.Manager,
	sel selector.Selector,
	upd updater.Updater,
) *App {
	return &App{
		config:   cfg,
		console:  console,
		depMgr:   depMgr,
		selector: sel,
		updater:  upd,
	}
}

// Run executes the main application logic
func (a *App) Run() error {
	a.console.Header()

	// Debug: Print configuration
	if a.config.Verbose {
		if a.config.List {
			a.console.Debug("List mode enabled")
		}
		if a.config.Selective {
			a.console.Debug("Selective mode enabled")
		}
		if a.config.Interactive {
			a.console.Debug("Interactive mode enabled")
		}
		if a.config.All {
			a.console.Debug("All dependencies mode enabled")
		}
	}

	// Get only updatable dependencies
	allUpdatableDeps, err := a.depMgr.GetUpdatableDependencies()
	if err != nil {
		return err
	}

	if len(allUpdatableDeps) == 0 {
		a.console.Info("All dependencies are up to date! 🎉")
		return nil
	}

	// Filter dependencies based on configuration (direct vs all)
	filteredDeps := a.depMgr.FilterDependencies(allUpdatableDeps, a.config.ShouldIncludeIndirect())
	if len(filteredDeps) == 0 {
		if a.config.ShouldIncludeIndirect() {
			a.console.Info("All dependencies are up to date! 🎉")
		} else {
			a.console.Info("All direct dependencies are up to date! 🎉")
			if len(allUpdatableDeps) > 0 {
				indirectCount := len(allUpdatableDeps) - len(filteredDeps)
				a.console.Info("(%d indirect dependencies have updates available, use --all to include them)", indirectCount)
			}
		}
		return nil
	}

	// Select dependencies to update
	a.console.Debug("Selecting dependencies to update...")
	selectedDeps, err := a.selectDependencies(filteredDeps)
	if err != nil {
		return fmt.Errorf("dependency selection failed: %w", err)
	}

	a.console.Debug("Selected %d dependencies for update", len(selectedDeps))

	if len(selectedDeps) == 0 {
		a.console.Info("No dependencies selected for update")
		return nil
	}

	// Handle List mode
	if a.config.List {
		return nil
	}

	// Confirm update if in interactive mode (but not selective, as that already confirms)
	if a.config.Interactive && !a.config.Selective {
		if !a.console.Confirm("Do you want to proceed with the update?") {
			a.console.Info("Update cancelled")
			return nil
		}
	}

	// Perform the update - handle failures gracefully
	return a.performUpdate(selectedDeps)
}

func (a *App) selectDependencies(deps []dependency.Dependency) ([]dependency.Dependency, error) {
	if !a.config.Selective {
		// Non-selective mode: show dependencies that will be updated and return all
		typeStr := "direct"
		if a.config.All {
			typeStr = "all"
		}
		title := fmt.Sprintf("Found %d %s dependencies with available updates:", len(deps), typeStr)
		a.console.PrintDependencies(deps, title)
		return deps, nil
	}

	// Selective mode: use interactive selection
	result := a.selector.Select(deps, a.config.ShouldIncludeIndirect())
	if result.Error != nil {
		return nil, result.Error
	}

	if result.Cancelled {
		return []dependency.Dependency{}, nil
	}

	return result.Selected, nil
}

func (a *App) performUpdate(deps []dependency.Dependency) error {
	a.console.Info("Updating dependencies...")

	// Update dependencies with progress reporting
	result := a.updateWithProgress(deps)

	// Report results
	a.console.PrintUpdateResult(len(result.Updated), len(deps), len(result.Failed) > 0)

	// Show individual errors if any - but don't fail the whole process
	for _, failure := range result.Failed {
		a.console.Error("Failed to update %s: %v", failure.Dependency.Path, failure.Error)
	}

	// Run go mod tidy - even if some updates failed
	if err := a.runModTidy(); err != nil {
		// Don't fail completely if mod tidy fails
		a.console.Warning("go mod tidy failed: %v", err)
	} else {
		a.console.Success("go mod tidy completed")
	}

	// Show final status
	if len(result.Updated) > 0 {
		a.console.Success("Dependency update completed!")
		if len(result.Failed) > 0 {
			a.console.Info("Successfully updated %d out of %d dependencies", len(result.Updated), len(deps))
		}
	} else if len(result.Failed) > 0 {
		a.console.Warning("No dependencies were successfully updated due to errors")
	}

	return nil // Don't fail the whole process for individual dependency issues
}

func (a *App) updateWithProgress(deps []dependency.Dependency) updater.UpdateResult {
	var allResults []updater.UpdateResult

	for i, dep := range deps {
		a.console.ProgressBar(i, len(deps), dep.Path)

		// Update individual dependency - errors are captured in result
		singleResult := a.updater.UpdateDependencies([]dependency.Dependency{dep}, a.config.Verbose)
		allResults = append(allResults, singleResult)

		a.console.ProgressBar(i+1, len(deps), dep.Path)
	}

	finalResult := updater.UpdateResult{
		Updated: make([]dependency.Dependency, 0),
		Failed:  make([]updater.UpdateError, 0),
		Success: true,
	}

	for _, result := range allResults {
		finalResult.Updated = append(finalResult.Updated, result.Updated...)
		finalResult.Failed = append(finalResult.Failed, result.Failed...)
		if !result.Success {
			finalResult.Success = false
		}
	}

	return finalResult
}

func (a *App) runModTidy() error {
	a.console.Info("Running go mod tidy...")
	return a.updater.RunModTidy(a.config.Verbose)
}
