package cmd

import "github.com/spf13/cobra"

var sharesCmd = &cobra.Command{
	Use:     "shares",
	Aliases: []string{"share", "sh"},
	Short:   "Commands to manage TOTP shares",
}

func init() {
	sharesCmd.AddCommand(
		sharesSplitCmd,
	)
}
