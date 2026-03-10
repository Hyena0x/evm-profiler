package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config module
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage evm-profiler configuration",
	Long:  `View or set configuration values like RPC URLs and API keys.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long: `Set configuration keys and save them to the persistent ~/.evm-profiler.yaml file.
For example:
evm-profiler config set --rpc "https://eth-mainnet.xyz" --etherscan-key "ABCDE"
`,
	Run: func(cmd *cobra.Command, args []string) {
		changesMade := false

		rpc, _ := cmd.Flags().GetString("rpc")
		if rpc != "" {
			viper.Set("rpc", rpc)
			changesMade = true
		}

		etherscanKey, _ := cmd.Flags().GetString("etherscan-key")
		if etherscanKey != "" {
			viper.Set("etherscan-key", etherscanKey)
			changesMade = true
		}

		goplusKey, _ := cmd.Flags().GetString("goplus-key")
		if goplusKey != "" {
			viper.Set("goplus-key", goplusKey)
			changesMade = true
		}

		network, _ := cmd.Flags().GetString("network")
		if network != "" {
			viper.Set("network", network)
			changesMade = true
		}

		if changesMade {
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, ".evm-profiler.yaml")

			err := viper.WriteConfigAs(configPath)
			if err != nil {
				fmt.Printf("Error saving configuration to %s: %s\n", configPath, err)
				return
			}
			fmt.Printf("Configuration saved to %s\n", configPath)
		} else {
			fmt.Println("No configuration values provided to set.")
			cmd.Help()
		}
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long:  `Print out all currently active configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Println("----------------------")

		fmt.Printf("RPC URL:        %s\n", viper.GetString("rpc"))
		fmt.Printf("Etherscan Key:  %s\n", viper.GetString("etherscan-key"))
		fmt.Printf("GoPlus Key:     %s\n", viper.GetString("goplus-key"))
		fmt.Printf("Network:        %s\n", viper.GetString("network"))

		configFile := viper.ConfigFileUsed()
		if configFile != "" {
			fmt.Printf("\nLoaded from: %s\n", configFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configViewCmd)

	configSetCmd.Flags().String("rpc", "", "RPC Node URL")
	configSetCmd.Flags().String("etherscan-key", "", "Etherscan API Key")
	configSetCmd.Flags().String("goplus-key", "", "GoPlus API Key")
	configSetCmd.Flags().String("network", "", "Network name (e.g., mainnet, arbitrum)")
}
