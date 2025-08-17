package config

// Config holds all configuration options for the application
type Config struct {
	DryRun      bool // Show what would be updated without making changes
	Interactive bool // Ask for confirmation before updating
	Verbose     bool // Show detailed output
	NoColor     bool // Disable colored output
	All         bool // Update indirect dependencies as well
	Selective   bool // Interactively select which dependencies to update
}

// ShouldIncludeIndirect returns true if indirect dependencies should be included
func (c *Config) ShouldIncludeIndirect() bool {
	return c.All
}

// IsInteractiveMode returns true if any interactive mode is enabled
func (c *Config) IsInteractiveMode() bool {
	return c.Interactive || c.Selective
}
