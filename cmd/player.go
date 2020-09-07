package cmd

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

func playerCmd(client *api.Client) *cobra.Command {
	c := &cobra.Command{
		Use:     "player",
		Aliases: []string{"play", "p"},
		Short:   "Players are those holding shares and helping a dealer retrieve a secret",
	}

	c.AddCommand(
		playerStartCmd(client),
	)

	return c
}
