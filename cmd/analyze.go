package cmd

import (
	"context"
	"fmt"
	"time"

	"evm-profiler/internal/analyzer"
	"evm-profiler/internal/fetcher"
	"evm-profiler/internal/printer"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var analyzeCmd = &cobra.Command{
	Use:     "analyze [address]",
	Aliases: []string{"scan"},
	Short:   "Generate a profile report for a specific EVM address",
	Long: `Analyze concurrently issues network requests to gather RPC, 
Etherscan, and GoPlus data, aggregating them into a terminal profile report.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		rpcURL := viper.GetString("rpc")
		if rpcURL == "" {
			return fmt.Errorf("RPC URL is required. Please provide it via --rpc flag or run `evm-profiler config set --rpc <url>`")
		}

		etherscanKey := viper.GetString("etherscan-key")
		if etherscanKey == "" {
			color.Yellow("Warning: Etherscan API key is empty. Analytics like tx count and funder may fail rate limits.")
		}

		// Initialize fetchers
		rpcFetcher, err := fetcher.NewRPCFetcher(rpcURL)
		if err != nil {
			return fmt.Errorf("failed to init RPC: %w", err)
		}

		etherscanFetcher := fetcher.NewEtherscanFetcher(etherscanKey)
		goplusFetcher := fetcher.NewGoPlusFetcher()

		analyzerSvc := analyzer.NewAnalyzer(rpcFetcher, etherscanFetcher, goplusFetcher)

		profile, err := analyzerSvc.Analyze(ctx, address)
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		printer.PrintProfile(profile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().String("rpc", "", "Ethereum RPC Node URL")
	analyzeCmd.Flags().String("etherscan-key", "", "Etherscan API Key")
	analyzeCmd.Flags().String("goplus-key", "", "GoPlus API Key")
	analyzeCmd.Flags().String("network", "mainnet", "EVM Network")

	viper.BindPFlag("rpc", analyzeCmd.Flags().Lookup("rpc"))
	viper.BindPFlag("etherscan-key", analyzeCmd.Flags().Lookup("etherscan-key"))
	viper.BindPFlag("goplus-key", analyzeCmd.Flags().Lookup("goplus-key"))
	viper.BindPFlag("network", analyzeCmd.Flags().Lookup("network"))
}
