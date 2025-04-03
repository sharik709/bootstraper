package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/sharik709/bootstraper/providers"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available frameworks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available frameworks:")
		fmt.Println("---------------------")
		
		for _, provider := range providers.List() {
			fmt.Printf("%-10s - %s\n", provider.Name(), provider.Description())
		}
	},
}

// cmd/project.go - Project command (alternative syntax)
package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/username/bootstraper/providers"
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
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if f.Name != framework {
				options[f.Name] = f.Value.String()
			}
		})
		
		// Bootstrap the project
		return provider.Bootstrap(projectName, options)
	},
}

func init() {
	// Add framework flags
	for _, provider := range providers.List() {
		projectCmd.Flags().Bool(provider.Name(), false, fmt.Sprintf("Create a %s project", provider.Name()))
		
		// Add framework-specific options to the new command
		for option, description := range provider.AvailableOptions() {
			newCmd.Flags().String(option, "", description)
		}
	}
	
	rootCmd.AddCommand(projectCmd)
}
