package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multifarmer",
	Short: "MFA for big farms",
}

func init() {
	rootCmd.AddCommand(
		serverCmd,
		dealerCmd,
		playerCmd,
		sharesCmd,
	)
}

// Execute is used by main as entrypoint
func Execute() error {
	return rootCmd.Execute()
}
