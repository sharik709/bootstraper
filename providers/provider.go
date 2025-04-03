package providers

import (
	"errors"
	"sort"
)

// Provider defines the interface for framework providers
type Provider interface {
	// Name returns the name of the framework provider
	Name() string

	// Description returns a brief description of the framework
	Description() string

	// Bootstrap initializes a new project with the given name and options
	Bootstrap(projectName string, options map[string]string) error

	// AvailableOptions returns a map of available options for the provider
	AvailableOptions() map[string]string

	// CheckDependencies verifies that all required dependencies are installed
	// Returns nil if all dependencies are met, or error with instructions
	CheckDependencies() error

	// SupportedVersions returns a list of supported framework versions
	// If empty, uses the latest version by default
	SupportedVersions() []string
}

// Registry keeps track of all registered providers
var Registry = make(map[string]Provider)

// Register adds a provider to the registry
func Register(p Provider) {
	Registry[p.Name()] = p
}

// Get returns a provider by name
func Get(name string) (Provider, error) {
	provider, exists := Registry[name]
	if !exists {
		return nil, errors.New("provider not found: " + name)
	}
	return provider, nil
}

// List returns all registered providers
func List() []Provider {
	providers := make([]Provider, 0, len(Registry))
	for _, provider := range Registry {
		providers = append(providers, provider)
	}

	// Sort by name for consistent output
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name() < providers[j].Name()
	})

	return providers
}
