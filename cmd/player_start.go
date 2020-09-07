package cmd

import (
	"fmt"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
)

func playerStartCmd(client *api.Client) *cobra.Command {

	run := func(cmd *cobra.Command, args []string) error {
		rawResp, _ := client.CreateSession(cmd.Context(), api.CreateSessionJSONRequestBody{
			FirstShare: "",
			Ttl:        nil,
		})

		resp, _ := api.ParseCreateSessionResponse(rawResp)

		fmt.Print(resp.JSON200.Shares)

		return nil
	}

	return &cobra.Command{
		Use:   "start",
		Short: "Start a new session",
		RunE:  run,
	}
}
