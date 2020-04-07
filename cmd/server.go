package cmd

import (
	"github.com/giorgioazzinnaro/multi-farmer-authentication/server"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve", "s"},
	Short:   "Start the server",

	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	e := echo.New()

	server.Handle(e)

	e.Logger.Fatal(e.Start(":8081"))

	return nil
}
