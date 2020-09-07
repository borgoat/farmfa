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

	c.PersistentFlags().StringP("address", "a", "http://localhost:8080", "The endpoint to the API")
	_ = viper.BindPFlag("address", c.PersistentFlags().Lookup("address"))
	_ = viper.BindEnv("address", "FARMFA_ADDRESS")

	viper.ReadInConfig()

	defaultClient, err := api.NewClient(viper.GetString("address"))
	if err != nil {
		panic(err)
	}

	c.AddCommand(
		playerCmd(defaultClient),
		dealerCmd(defaultClient),
		sharesCmd(defaultClient),

		serverCmd(&ServerConfig{}),
	)

	return c
}

// Execute is used by main as entrypoint
func Execute() error {
	return rootCmd().Execute()
}
