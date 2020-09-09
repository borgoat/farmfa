package cmd

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

func dealerCmd(client *api.Client) *cobra.Command {
	c := &cobra.Command{
		Use:     "dealer",
		Aliases: []string{"deal", "d"},
		Short:   "Dealers need help from players to retrieve secrets, they initiate sessions",
	}

	c.AddCommand(
		dealerStartCmd(client),
	)

	return c
}
