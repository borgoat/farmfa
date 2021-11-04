package cmd

import (
	"fmt"

	"github.com/borgoat/farmfa/api"
	"github.com/borgoat/farmfa/server"
	"github.com/borgoat/farmfa/session"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd(v *viper.Viper) *cobra.Command {

	run := func(cmd *cobra.Command, args []string) error {
		e := echo.New()

		apiObj, err := api.GetSwagger()
		if err != nil {
			return fmt.Errorf("error loading OpenAPI spec: %w", err)
		}
		e.Use(middleware.OapiRequestValidator(apiObj))

		store := session.NewInMemoryStore()
		oracle := session.NewOracle(store)

		s := server.New(oracle)

		api.RegisterHandlers(e, s)

		e.Logger.Fatal(e.Start(v.GetString("bind-address")))

		return nil
	}

	c := &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve", "s", "oracle"},
		Short:   "Start the server acting as the farMFA oracle",

		Long: `The oracle is the entity that reconstructs Tocs into TOTP secrets, and generates one-time passwords.
Also called the prover, as defined in [RFC6238].`,

		RunE: run,
	}

	c.Flags().String("bind-address", "localhost:8080", "The address to bind the server")
	_ = v.BindPFlags(c.Flags())
	_ = v.BindEnv("bind-address", "FARMFA_BIND_ADDRESS")
	_ = v.ReadInConfig()

	return c
}
