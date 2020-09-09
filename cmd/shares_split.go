package cmd

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func sharesSplitCmd(v *viper.Viper, client *api.Client) *cobra.Command {

	run := func(cmd *cobra.Command, args []string) error {
		secret := v.GetString("totp-secret")
		if secret == "" {
			r, _ := readline.New("Paste a TOTP secret you want to split > ")
			secret, _ = r.Readline()
		}

		// TODO Input users you want to encrypt for, then threshold

		rawResp, _ := client.CreateShares(cmd.Context(), api.CreateSharesJSONRequestBody{
			EncryptionKeys: nil,
			Shares:         5,
			Threshold:      3,
			TotpSecretKey:  secret,
		})

		resp, _ := api.ParseCreateSharesResponse(rawResp)

		fmt.Print(resp.JSON200.Shares)

		return nil
	}

	c := &cobra.Command{
		Use:   "split",
		Short: "Split a new TOTP",
		RunE:  run,
	}

	c.Flags().Bool("disable-encryption", false, "Add this flag to return unencrypted shares")

	c.Flags().String("totp-secret", "", "The secret to be split, if not provided, it will be requested interactively")
	_ = v.BindPFlags(c.Flags())
	_ = v.BindEnv("totp-secret", "FARMFA_TOTP_SECRET")
	_ = v.ReadInConfig()

	return c
}
