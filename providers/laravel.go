package providers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sharik709/bootstraper/util"
)

// LaravelProvider provider
type LaravelProvider struct{}

func init() {
	Register(&LaravelProvider{})
}

func (p *LaravelProvider) Name() string {
	return "laravel"
}

func (p *LaravelProvider) Description() string {
	return "Laravel - PHP web application framework"
}

func (p *LaravelProvider) Bootstrap(projectName string, options map[string]string) error {
	// Check if Laravel installer or Composer is available
	if util.CommandExists("laravel") {
		// Use Laravel installer
		cmd := exec.Command("laravel", "new", projectName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		fmt.Printf("Creating Laravel project: %s\n", projectName)
		return cmd.Run()
	} else if util.CommandExists("composer") {
		// Fall back to composer
		cmd := exec.Command("composer", "create-project", "laravel/laravel", projectName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		fmt.Printf("Creating Laravel project: %s\n", projectName)
		return cmd.Run()
	} else {
		return fmt.Errorf("neither Laravel installer nor Composer found. Please install one of them")
	}
}

func (p *LaravelProvider) AvailableOptions() map[string]string {
	return map[string]string{
		"git": "Initialize a Git repository (true/false)",
	}
}
