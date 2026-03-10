package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of evm-profiler",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("evm-profiler v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
