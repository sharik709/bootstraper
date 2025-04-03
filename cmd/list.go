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
