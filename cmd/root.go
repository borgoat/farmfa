package cmd

import (
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

	c.AddCommand(
		serverCmd(v),
		dealerCmd(v),
	)

	return c
}

// Execute is used by main as entrypoint
func Execute() error {
	return rootCmd().Execute()
}
