package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sharik709/bootstraper/util"
)

type ProviderDefinition struct {
	ProviderName string            `json:"name"`
	ProviderDesc string            `json:"description"`
	Command      string            `json:"command"`
	CommandArgs  []string          `json:"args"`
	DependsOn    []string          `json:"dependencies"`
	Options      map[string]string `json:"options"`
	Versions     []string          `json:"versions,omitempty"`
}

type ProviderRegistry struct {
	Providers []ProviderDefinition `json:"providers"`
	UpdatedAt string               `json:"updated_at"`
}

func (p *ProviderDefinition) Name() string {
	return p.ProviderName
}

func (p *ProviderDefinition) Description() string {
	return p.ProviderDesc
}

func (p *ProviderDefinition) Bootstrap(projectName string, options map[string]string) error {
	for _, dep := range p.DependsOn {
		if !util.CommandExists(dep) {
			return fmt.Errorf("dependency not found: %s", dep)
		}
	}

	args := make([]string, len(p.CommandArgs))
	copy(args, p.CommandArgs)

	// Replace project name placeholder
	for i, arg := range args {
		if arg == "{project-name}" {
			args[i] = projectName
		}
	}

	// Handle version placeholders
	version, hasVersion := options["version"]
	for i, arg := range args {
		if strings.Contains(arg, "{version}") {
			if hasVersion && version != "" {
				args[i] = strings.ReplaceAll(arg, "{version}", version)
			} else {
				// Use "latest" if no version specified or remove version tag altogether
				if strings.Contains(arg, "@{version}") {
					// For patterns like package@{version}, use latest
					args[i] = strings.ReplaceAll(arg, "@{version}", "@latest")
				} else {
					// For other cases, just remove the {version} placeholder
					args[i] = strings.ReplaceAll(arg, "{version}", "")
				}
			}
		}
	}

	// Handle module placeholder (specific to Go)
	if module, ok := options["module"]; ok && module != "" {
		for i, arg := range args {
			if arg == "{module}" {
				args[i] = module
			}
		}
	} else {
		// If module is a placeholder but not provided, use a default based on project name
		for i, arg := range args {
			if arg == "{module}" {
				args[i] = "github.com/example/" + projectName
			}
		}
	}

	// Add command-line flags from options
	for flagName, flagValue := range options {
		if flagName == "version" || flagName == "module" {
			continue // Already handled
		}

		if flagValue == "true" {
			args = append(args, fmt.Sprintf("--%s", flagName))
		} else if flagValue != "false" && flagValue != "" {
			args = append(args, fmt.Sprintf("--%s=%s", flagName, flagValue))
		}
	}

	cmd := exec.Command(p.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Creating %s project: %s\n", p.ProviderName, projectName)
	return cmd.Run()
}

func (p *ProviderDefinition) AvailableOptions() map[string]string {
	return p.Options
}

func (p *ProviderDefinition) CheckDependencies() error {
	for _, dep := range p.DependsOn {
		if !util.CommandExists(dep) {
			return fmt.Errorf("dependency not found: %s", dep)
		}
	}
	return nil
}

func (p *ProviderDefinition) SupportedVersions() []string {
	return p.Versions
}

func loadProviders() ([]Provider, error) {
	registryPath := "providers/registry.json"

	data, err := os.ReadFile(registryPath)
	if err != nil {
		// Try with absolute path from executable
		executablePath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}

		executableDir := filepath.Dir(executablePath)
		registryPath = filepath.Join(executableDir, "providers", "registry.json")

		data, err = os.ReadFile(registryPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read provider registry: %v", err)
		}
	}

	var registry ProviderRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse provider registry: %v", err)
	}

	providers := make([]Provider, 0, len(registry.Providers))
	for _, providerDef := range registry.Providers {
		provider := providerDef // Create local copy to avoid referencing loop variable
		providers = append(providers, &provider)
	}

	return providers, nil
}

func init() {
	// Load providers from registry.json
	jsonProviders, err := loadProviders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load providers from registry.json: %v\n", err)
		return
	}

	for _, provider := range jsonProviders {
		Register(provider)
	}
}
