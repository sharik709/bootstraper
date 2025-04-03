package cmd

import (
	"fmt"

	"github.com/sharik709/bootstraper/providers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		cmd.Flags().Visit(func(f *pflag.Flag) {
			options[f.Name] = f.Value.String()
		})

		// Bootstrap the project
		return provider.Bootstrap(projectName, options)
	},
}

func init() {
	// Dynamic flags will be added when providers are registered
}
