package providers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/username/bootstraper/util"
)

// VueProvider provider
type VueProvider struct{}

func init() {
	Register(&VueProvider{})
}

func (p *VueProvider) Name() string {
	return "vue"
}

func (p *VueProvider) Description() string {
	return "Vue.js - Progressive JavaScript framework"
}

func (p *VueProvider) Bootstrap(projectName string, options map[string]string) error {
	// Check if npm is available
	if !util.CommandExists("npm") {
		return fmt.Errorf("npm command not found. Please install Node.js and npm")
	}

	// Build the command: npm create vue@latest projectName
	cmd := exec.Command("npm", "create", "vue@latest", projectName)

	// Set the output to the current terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	fmt.Printf("Creating Vue.js project: %s\n", projectName)
	return cmd.Run()
}

func (p *VueProvider) AvailableOptions() map[string]string {
	return map[string]string{
		// Vue CLI handles options interactively
	}
}
