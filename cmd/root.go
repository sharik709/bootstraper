package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bt",
	Short: "Bootstraper - A universal project bootstrapping tool",
	Long: `Bootstraper (bt) is a unified CLI tool designed to simplify and standardize
the process of initializing new development projects across multiple
frameworks, languages, and platforms.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
}
