package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sharik709/bootstraper/util"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage custom project templates",
	Long: `Manage custom project templates.

  Templates allow you to create projects from custom sources like git repositories
  or local directories, with optional pre-defined values.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		if len(config.Templates) == 0 {
			fmt.Println("No templates configured. Use 'bt template add' to add a template.")
			return nil
		}

		fmt.Println("Available templates:")
		fmt.Println("--------------------")
		for name, template := range config.Templates {
			fmt.Printf("%-20s - %s\n", name, template.Description)
			fmt.Printf("                      Source: %s\n", template.Source)
			if len(template.Tags) > 0 {
				fmt.Printf("                      Tags: %s\n", strings.Join(template.Tags, ", "))
			}
			fmt.Println()
		}

		return nil
	},
}

var templateAddCmd = &cobra.Command{
	Use:   "add [name] [source]",
	Short: "Add a new template",
	Long: `Add a new template from a source.

  The source can be a git repository, local directory, or archive file.
  For example:
    bt template add my-nextjs github:username/my-nextjs-template
    bt template add my-flask /path/to/local/template
  `,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		source := args[1]

		// Load configuration
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Get other flags
		description, _ := cmd.Flags().GetString("description")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// Create template
		if config.Templates == nil {
			config.Templates = make(map[string]util.Template)
		}

		config.Templates[name] = util.Template{
			Source:      source,
			Description: description,
			Tags:        tags,
		}

		// Save configuration
		if err := util.SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Template '%s' added successfully.\n", name)
		return nil
	},
}

var templateRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Load configuration
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Check if template exists
		if _, ok := config.Templates[name]; !ok {
			return fmt.Errorf("template '%s' not found", name)
		}

		// Remove template
		delete(config.Templates, name)

		// Save configuration
		if err := util.SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Template '%s' removed successfully.\n", name)
		return nil
	},
}

var templateUseCmd = &cobra.Command{
	Use:   "use [template] [project-name]",
	Short: "Create a project from a template",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		projectName := args[1]

		// Load configuration
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Check if template exists
		template, ok := config.Templates[templateName]
		if !ok {
			return fmt.Errorf("template '%s' not found", templateName)
		}

		// Create project directory
		if err := os.MkdirAll(projectName, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}

		// Clone/copy from source
		if strings.HasPrefix(template.Source, "github:") {
			// Handle GitHub repository
			repo := strings.TrimPrefix(template.Source, "github:")
			gitURL := fmt.Sprintf("https://github.com/%s.git", repo)

			cmd := exec.Command("git", "clone", gitURL, projectName)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("Creating project from GitHub template: %s\n", repo)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to clone repository: %v", err)
			}

			// Remove .git directory
			if err := os.RemoveAll(filepath.Join(projectName, ".git")); err != nil {
				fmt.Printf("Warning: failed to remove .git directory: %v\n", err)
			}
		} else if strings.HasPrefix(template.Source, "http://") || strings.HasPrefix(template.Source, "https://") {
			// Handle remote URL
			return fmt.Errorf("direct URL download not implemented yet")
		} else {
			// Handle local directory
			// For simplicity, we'll just copy the directory recursively
			return fmt.Errorf("local directory copy not implemented yet")
		}

		fmt.Printf("Project '%s' created from template '%s'.\n", projectName, templateName)
		return nil
	},
}

func init() {
	// Configure template add command
	templateAddCmd.Flags().String("description", "", "Description of the template")
	templateAddCmd.Flags().StringSlice("tags", []string{}, "Tags for categorizing the template")

	// Add subcommands
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateRemoveCmd)
	templateCmd.AddCommand(templateUseCmd)

	// Add to root command
	rootCmd.AddCommand(templateCmd)
}
