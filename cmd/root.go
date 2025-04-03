package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	// Version is the current version of the CLI
	Version = "0.2.0"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:     "bt",
	Version: Version,
	Short:   "Bootstraper - A universal project bootstrapping tool",
	Long: `Bootstraper (bt) is a unified CLI tool designed to simplify and standardize
the process of initializing new development projects across multiple
frameworks, languages, and platforms.

📦 Usage Examples:
  bt new next my-nextjs-app           Create a Next.js project
  bt new vue my-vue-app               Create a Vue.js project
  bt list                             List available frameworks
  bt config set defaults.next.typescript true  Set default options`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, display help
		cmd.Help()
	},
}

func Execute() error {
	// Add custom formatting for version output
	rootCmd.SetVersionTemplate(`Bootstraper v{{.Version}}
  https://github.com/sharik709/bootstraper
`)

	// Handle any errors during command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Register commands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(projectCmd)
	// Additional commands added in their respective files:
	// - configCmd
	// - templateCmd
}
