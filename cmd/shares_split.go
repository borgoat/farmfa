package cmd

import (
	"fmt"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

func sharesSplitCmd(client *api.Client) *cobra.Command {

	run := func(cmd *cobra.Command, args []string) error {
		rawResp, _ := client.CreateShares(cmd.Context(), api.CreateSharesJSONRequestBody{
			EncryptionKeys: nil,
			Shares:         5,
			Threshold:      3,
			TotpSecretKey:  "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
		})

		resp, _ := api.ParseCreateSharesResponse(rawResp)

		fmt.Print(resp.JSON200.Shares)

		return nil
	}

	return &cobra.Command{
		Use:   "split",
		Short: "Split a new TOTP",
		RunE:  run,
	}
}
