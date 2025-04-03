package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sharik709/bootstraper/util"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage bootstraper configuration",
	Long: `Manage bootstraper configuration settings.
For example:
  bt config get
  bt config set defaults.next.typescript true
  bt config reset`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Long: `Get configuration value by key.
For example:
  bt config get defaults.next.typescript
  bt config get telemetry
  bt config get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		if len(args) == 0 {
			// Show entire config
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %v", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Parse the key path
		path := strings.Split(args[0], ".")
		var result interface{} = config

		// Navigate through the config
		for _, key := range path {
			switch v := result.(type) {
			case map[string]interface{}:
				result = v[key]
			case map[string]map[string]interface{}:
				result = v[key]
			case util.Config:
				switch key {
				case "defaults":
					result = v.Defaults
				case "templates":
					result = v.Templates
				case "telemetry":
					result = v.Telemetry
				case "cacheDir":
					result = v.CacheDir
				case "projectDir":
					result = v.ProjectDir
				default:
					return fmt.Errorf("key not found: %s", args[0])
				}
			default:
				return fmt.Errorf("key not found: %s", args[0])
			}
		}

		// Print the result
		if result == nil {
			fmt.Println("null")
		} else {
			switch v := result.(type) {
			case map[string]interface{}, map[string]map[string]interface{}, map[string]util.Template:
				data, _ := json.MarshalIndent(v, "", "  ")
				fmt.Println(string(data))
			default:
				fmt.Println(v)
			}
		}

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set configuration value",
	Long: `Set configuration value by key.
For example:
  bt config set defaults.next.typescript true
  bt config set telemetry false
  bt config set projectDir ~/Projects`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := util.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		key := args[0]
		value := args[1]

		// Handle special case for default provider options
		if strings.HasPrefix(key, "defaults.") {
			parts := strings.Split(key, ".")
			if len(parts) < 3 {
				return fmt.Errorf("invalid key format for defaults: %s", key)
			}

			providerName := parts[1]
			optionName := parts[2]

			defaults, err := util.GetDefaultsForProvider(providerName)
			if err != nil {
				return err
			}

			// Convert value to appropriate type
			var typedValue interface{}
			switch value {
			case "true", "yes", "1":
				typedValue = true
			case "false", "no", "0":
				typedValue = false
			default:
				// Try to parse as number
				var numValue float64
				if err := json.Unmarshal([]byte(value), &numValue); err == nil {
					typedValue = numValue
				} else {
					// Treat as string
					typedValue = value
				}
			}

			defaults[optionName] = typedValue
			if err := util.SetDefaultsForProvider(providerName, defaults); err != nil {
				return err
			}

			fmt.Printf("Set %s to %v\n", key, typedValue)
			return nil
		}

		// Handle other config settings
		switch key {
		case "telemetry":
			config.Telemetry = value == "true" || value == "yes" || value == "1"
		case "cacheDir":
			config.CacheDir = os.ExpandEnv(value)
		case "projectDir":
			config.ProjectDir = os.ExpandEnv(value)
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		if err := util.SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Set %s to %s\n", key, value)
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := util.DefaultConfig()
		if err := util.SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Println("Configuration reset to defaults")
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
