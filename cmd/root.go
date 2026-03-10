package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "evm-profiler",
	Short: "A stateless, high-performance Golang CLI tool to profile EVM addresses.",
	Long: `evm-profiler concurrently fetches data from RPC, Etherscan
and GoPlus to produce a comprehensive terminal report of an EVM address.`,
}

// Execute adds all child commands to the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".evm-profiler" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".evm-profiler")

	// Set defaults
	viper.SetDefault("network", "mainnet")

	// Read in environment variables that match
	viper.SetEnvPrefix("EVM_PROFILER")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Only log for debugging, silently load otherwise
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// If no config file exists, create an empty one to avoid errors later on save
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configPath := filepath.Join(home, ".evm-profiler.yaml")
			os.WriteFile(configPath, []byte(""), 0644)
		}
	}
}
