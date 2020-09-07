package cmd

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

func sharesCmd(client *api.Client) *cobra.Command {
	c := &cobra.Command{
		Use:     "shares",
		Aliases: []string{"share", "sh"},
		Short:   "Commands to manage TOTP shares",
	}

	c.AddCommand(
		sharesSplitCmd(client),
	)

	return c
}
