package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldIncludeIndirect(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "all flag enabled",
			config:   Config{All: true},
			expected: true,
		},
		{
			name:     "all flag disabled",
			config:   Config{All: false},
			expected: false,
		},
		{
			name:     "default config",
			config:   Config{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ShouldIncludeIndirect()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInteractiveMode(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "interactive flag enabled",
			config:   Config{Interactive: true},
			expected: true,
		},
		{
			name:     "selective flag enabled",
			config:   Config{Selective: true},
			expected: true,
		},
		{
			name:     "both flags enabled",
			config:   Config{Interactive: true, Selective: true},
			expected: true,
		},
		{
			name:     "no interactive flags",
			config:   Config{Interactive: false, Selective: false},
			expected: false,
		},
		{
			name:     "default config",
			config:   Config{},
			expected: false,
		},
		{
			name: "other flags don't affect interactive mode",
			config: Config{
				DryRun:  true,
				Verbose: true,
				NoColor: true,
				All:     true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsInteractiveMode()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultValues(t *testing.T) {
	config := Config{}

	assert.False(t, config.DryRun)
	assert.False(t, config.Interactive)
	assert.False(t, config.Verbose)
	assert.False(t, config.NoColor)
	assert.False(t, config.All)
	assert.False(t, config.Selective)
}

func TestAllFieldsCombination(t *testing.T) {
	config := Config{
		DryRun:      true,
		Interactive: true,
		Verbose:     true,
		NoColor:     true,
		All:         true,
		Selective:   true,
	}

	assert.True(t, config.DryRun)
	assert.True(t, config.Interactive)
	assert.True(t, config.Verbose)
	assert.True(t, config.NoColor)
	assert.True(t, config.All)
	assert.True(t, config.Selective)
	assert.True(t, config.ShouldIncludeIndirect())
	assert.True(t, config.IsInteractiveMode())
}
