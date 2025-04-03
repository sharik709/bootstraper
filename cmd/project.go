package cmd

import (
	"fmt"

	"github.com/sharik709/bootstraper/providers"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [project-name] --[framework]",
	Short: "Create a new project with the specified framework",
	Long: `Create a new project using the specified framework.
For example:
  bt project my-app --next
  bt project my-app --vue
  bt project my-app --laravel`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Determine which framework flag was used
		var framework string
		for _, provider := range providers.List() {
			flagValue, err := cmd.Flags().GetBool(provider.Name())
			if err == nil && flagValue {
				framework = provider.Name()
				break
			}
		}

		if framework == "" {
			return fmt.Errorf("no framework specified, use --[framework] flag")
		}

		// Get the provider
		provider, err := providers.Get(framework)
		if err != nil {
			return err
		}

		// Collect options from flags (exclude framework flags)
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

	// Add framework flags
	for _, provider := range providers.List() {
		projectCmd.Flags().Bool(provider.Name(), false, fmt.Sprintf("Create a %s project", provider.Name()))

		// Add framework-specific options to the project command
		for option, description := range provider.AvailableOptions() {
			// Skip if flag already exists
			if _, exists := addedFlags[option]; !exists {
				projectCmd.Flags().String(option, "", description)
				addedFlags[option] = true
			}
		}
	}

	rootCmd.AddCommand(projectCmd)
}
