package cmd

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/server"
	"github.com/giorgioazzinnaro/farmfa/sessions"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

type ServerConfig struct {
	BindAddress string `json:"bind_address"`
	LogLevel    string `json:"log_level"`
}

func serverCmd(cfg *ServerConfig) *cobra.Command {

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

		e.Logger.Fatal(e.Start(cfg.BindAddress))

		return nil
	}

	return &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve", "s"},
		Short:   "Start the server",

		RunE: run,
	}
}
