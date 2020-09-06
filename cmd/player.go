package cmd

import (
	"github.com/spf13/cobra"
)

var playerCmd = &cobra.Command{
	Use:     "player",
	Aliases: []string{"play", "p"},
	Short:   "Players are those holding shares and helping a dealer retrieve a secret",
}
