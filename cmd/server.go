package cmd

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/server"
	"github.com/giorgioazzinnaro/farmfa/sessions"
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

		sessionManager := sessions.NewInMemory()

		s := server.New(sessionManager)

		api.RegisterHandlers(e, s)

		e.Logger.Fatal(e.Start(v.GetString("bind-address")))

		return nil
	}

	c := &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve", "s"},
		Short:   "Start the server",

		RunE: run,
	}

	c.Flags().String("bind-address", "localhost:8080", "The address to bind the server")
	_ = v.BindPFlags(c.Flags())
	_ = v.BindEnv("bind-address", "FARMFA_BIND_ADDRESS")
	_ = v.ReadInConfig()

	return c
}
