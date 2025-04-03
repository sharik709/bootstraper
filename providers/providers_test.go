package providers

import (
	"testing"
)

func TestProviderRegistry(t *testing.T) {
	t.Run("Register and Get provider", func(t *testing.T) {
		Registry = make(map[string]Provider) // Reset registry for test

		mockProvider := &MockProvider{
			NameValue:        "mockprovider",
			DescriptionValue: "Mock Provider",
		}

		Register(mockProvider)

		if len(Registry) != 1 {
			t.Errorf("Expected registry to have 1 provider, got %d", len(Registry))
		}

		provider, err := Get("mockprovider")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if provider != mockProvider {
			t.Errorf("Expected to get the same provider that was registered")
		}
	})

	t.Run("Get non-existent provider", func(t *testing.T) {
		Registry = make(map[string]Provider) // Reset registry for test

		_, err := Get("nonexistent")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("List providers", func(t *testing.T) {
		Registry = make(map[string]Provider) // Reset registry for test

		mockProvider1 := &MockProvider{
			NameValue:        "provider1",
			DescriptionValue: "Provider 1",
		}

		mockProvider2 := &MockProvider{
			NameValue:        "provider2",
			DescriptionValue: "Provider 2",
		}

		Register(mockProvider1)
		Register(mockProvider2)

		providers := List()

		if len(providers) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(providers))
		}
	})
}

func TestRemoteProviders(t *testing.T) {
	t.Run("Remote provider adapter", func(t *testing.T) {
		adapter := &ProviderAdapter{
			NameFunc:        func() string { return "remote" },
			DescriptionFunc: func() string { return "Remote Provider" },
			AvailableOptionsFunc: func() map[string]string {
				return map[string]string{"option1": "Option 1"}
			},
			CheckDependenciesFunc: func() error { return nil },
			SupportedVersionsFunc: func() []string { return []string{"1.0", "2.0"} },
		}

		if adapter.Name() != "remote" {
			t.Errorf("Expected name 'remote', got '%s'", adapter.Name())
		}

		if adapter.Description() != "Remote Provider" {
			t.Errorf("Expected description 'Remote Provider', got '%s'", adapter.Description())
		}

		options := adapter.AvailableOptions()
		if len(options) != 1 || options["option1"] != "Option 1" {
			t.Errorf("Available options not as expected")
		}

		err := adapter.CheckDependencies()
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		versions := adapter.SupportedVersions()
		if len(versions) != 2 || versions[0] != "1.0" || versions[1] != "2.0" {
			t.Errorf("Supported versions not as expected")
		}
	})
}

// MockProvider is a mock implementation of the Provider interface for testing
type MockProvider struct {
	NameValue        string
	DescriptionValue string
	BootstrapError   error
	Options          map[string]string
	Versions         []string
}

func (p *MockProvider) Name() string {
	return p.NameValue
}

func (p *MockProvider) Description() string {
	return p.DescriptionValue
}

func (p *MockProvider) Bootstrap(projectName string, options map[string]string) error {
	return p.BootstrapError
}

func (p *MockProvider) AvailableOptions() map[string]string {
	return p.Options
}

func (p *MockProvider) CheckDependencies() error {
	return nil
}

func (p *MockProvider) SupportedVersions() []string {
	return p.Versions
}
