package providers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sharik709/bootstraper/util"
)

// NextJSProvider provider
type NextJSProvider struct{}

func init() {
	Register(&NextJSProvider{})
}

func (p *NextJSProvider) Name() string {
	return "next"
}

func (p *NextJSProvider) Description() string {
	return "Next.js - React framework with server-side rendering"
}

func (p *NextJSProvider) Bootstrap(projectName string, options map[string]string) error {
	// Check if npm/npx is available
	if !util.CommandExists("npx") {
		return fmt.Errorf("npx command not found. Please install Node.js and npm")
	}

	// Build the command: npx create-next-app@latest projectName
	cmd := exec.Command("npx", "create-next-app@latest", projectName)

	// Add typescript option if specified
	if ts, ok := options["typescript"]; ok && ts == "true" {
		cmd.Args = append(cmd.Args, "--typescript")
	}

	// Set the output to the current terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	fmt.Printf("Creating Next.js project: %s\n", projectName)
	return cmd.Run()
}

func (p *NextJSProvider) AvailableOptions() map[string]string {
	return map[string]string{
		"typescript": "Use TypeScript (true/false)",
	}
}
