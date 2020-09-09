package cmd

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func sharesCmd(v *viper.Viper, client *api.Client) *cobra.Command {
	c := &cobra.Command{
		Use:     "shares",
		Aliases: []string{"share", "sh"},
		Short:   "Commands to manage TOTP shares",
	}

	c.AddCommand(
		sharesSplitCmd(v, client),
	)

	return c
}
