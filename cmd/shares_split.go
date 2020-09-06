package cmd

import (
	"fmt"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

var sharesSplitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a new TOTP",

	RunE: runSharesSplit,
}

func runSharesSplit(cmd *cobra.Command, args []string) error {

	c, _ := api.NewClient("http://localhost:8081")

	rawResp, _ := c.CreateShares(cmd.Context(), api.CreateSharesJSONRequestBody{
		EncryptionKeys: nil,
		Shares:         5,
		Threshold:      3,
		TotpSecretKey:  "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	})

	resp, _ := api.ParseCreateSharesResponse(rawResp)

	fmt.Print(resp.JSON200.Shares)

	return nil
}
