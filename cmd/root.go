package cmd

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "farmfa",
		Short: "Far away MFA",
	}

	v := viper.New()

	c.PersistentFlags().StringP("address", "a", "http://localhost:8080", "The endpoint to the API")
	_ = v.BindPFlag("address", c.PersistentFlags().Lookup("address"))
	_ = v.BindEnv("address", "FARMFA_ADDRESS")

	_ = v.ReadInConfig()

	defaultClient, err := api.NewClient(v.GetString("address"))
	if err != nil {
		panic(err)
	}

	c.AddCommand(
		playerCmd(defaultClient),
		dealerCmd(defaultClient),
		sharesCmd(v, defaultClient),

		serverCmd(v),
	)

	return c
}

// Execute is used by main as entrypoint
func Execute() error {
	return rootCmd().Execute()
}
