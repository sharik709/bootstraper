package cmd

import (
	"fmt"

	"github.com/sharik709/bootstraper/providers"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [framework] [project-name]",
	Short: "Create a new project with the specified framework",
	Long: `Create a new project using the specified framework.
For example:
  bt new next my-app
  bt new vue my-app
  bt new laravel my-app
  bt new go my-app --module=github.com/username/my-app`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		framework := args[0]
		projectName := args[1]

		// Get the provider
		provider, err := providers.Get(framework)
		if err != nil {
			return fmt.Errorf("framework not supported: %s\nRun 'bt list' to see available frameworks", framework)
		}

		// Collect options from flags
		options := make(map[string]string)

		// Process flag values for the chosen framework
		for optName := range provider.AvailableOptions() {
			// Check if the flag exists and was set
			if flag := cmd.Flag(optName); flag != nil && flag.Changed {
				options[optName] = flag.Value.String()
			}
		}

		// Bootstrap the project
		return provider.Bootstrap(projectName, options)
	},
}

func init() {
	// Create a map to track which flags have been added to avoid duplicates
	addedFlags := make(map[string]bool)

	// Add available options for each provider as flags
	for _, provider := range providers.List() {
		for option, description := range provider.AvailableOptions() {
			// Skip if flag already exists
			if _, exists := addedFlags[option]; !exists {
				newCmd.Flags().String(option, "", description)
				addedFlags[option] = true
			}
		}
	}
}
