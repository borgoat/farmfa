package cmd

import (
	"github.com/spf13/cobra"
)

var dealerCmd = &cobra.Command{
	Use:     "dealer",
	Aliases: []string{"deal", "d"},
	Short:   "Dealers need help from players to retrieve secrets, they initiate sessions",
}
