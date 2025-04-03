package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the user configuration
type Config struct {
	Defaults   map[string]map[string]interface{} `json:"defaults"`
	Templates  map[string]Template               `json:"templates"`
	Telemetry  bool                              `json:"telemetry"`
	CacheDir   string                            `json:"cacheDir"`
	ProjectDir string                            `json:"projectDir"`
}

// Template represents a custom project template
type Template struct {
	Source      string   `json:"source"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Defaults:   make(map[string]map[string]interface{}),
		Templates:  make(map[string]Template),
		Telemetry:  true,
		CacheDir:   filepath.Join(homeDir, ".bootstraper", "cache"),
		ProjectDir: filepath.Join(homeDir, "Projects"),
	}
}

// LoadConfig loads the configuration from the given file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the given file
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	return filepath.Join(homeDir, ".bootstraperrc"), nil
}

// GetDefaultsForProvider returns the default options for a provider
func GetDefaultsForProvider(providerName string) (map[string]interface{}, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	if defaults, ok := config.Defaults[providerName]; ok {
		return defaults, nil
	}

	return make(map[string]interface{}), nil
}

// SetDefaultsForProvider sets the default options for a provider
func SetDefaultsForProvider(providerName string, defaults map[string]interface{}) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if config.Defaults == nil {
		config.Defaults = make(map[string]map[string]interface{})
	}

	config.Defaults[providerName] = defaults
	return SaveConfig(config)
}
