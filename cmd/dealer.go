package cmd

import (
	"filippo.io/age"
	"fmt"
	"github.com/borgoat/farmfa/deal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagNote       = "comment"
	FlagTOTPSecret = "totp-secret"
	FlagPlayers    = "players"
	FlagThreshold  = "threshold"
)

func dealerCmd(v *viper.Viper) *cobra.Command {
	run := func(cmd *cobra.Command, args []string) error {
		note := v.GetString(FlagNote)
		totpSecret := v.GetString(FlagTOTPSecret)
		playersInput := v.GetStringMapString(FlagPlayers)
		threshold := v.GetInt(FlagThreshold)

		var players = make([]*deal.Player, len(playersInput))

		i := 0
		for name, key := range playersInput {
			id, err := age.ParseX25519Recipient(key)
			if err != nil {
				return fmt.Errorf("the key provided for %s is invalid: %w", name, err)
			}
			player, err := deal.NewPlayer(name, deal.EncryptWithAge(id))
			if err != nil {
				return fmt.Errorf("failed to create player %s: %w", name, err)
			}

			players[i] = player
			i++
		}

		tocs, err := deal.CreateTocs(note, totpSecret, players, threshold)
		if err != nil {
			return fmt.Errorf("failed to generate Tocs: %w", err)
		}

		cmd.Printf("Tocs: %+v", tocs)

		return nil
	}

	c := &cobra.Command{
		Use:     "dealer",
		Aliases: []string{"deal", "d"},
		Short:   "Split a TOTP secret among multiple players",

		Long: `Create n Tocs from a TOTP secret. These Tocs are then shared with players.
Multiple players will then be able to get together to regenerate a TOTP from these Tocs.`,

		RunE: run,
	}

	c.Flags().StringP(FlagNote, "c", "", "A note for the Tocs")
	c.Flags().String(FlagTOTPSecret, "", "The TOTP Secret (preferably via env variable)")
	c.Flags().StringToStringP(
		FlagPlayers,
		"p",
		map[string]string{},
		"The players that will receive encrypted Tocs. Map of names -> Age public key.",
	)
	c.Flags().IntP(FlagThreshold, "t", 3, "The minimum number of Tocs to reconstruct the TOTP")
	_ = v.BindPFlags(c.Flags())
	_ = v.BindEnv(FlagNote, "FARMFA_COMMENT")
	_ = v.BindEnv(FlagTOTPSecret, "FARMFA_TOTP_SECRET")
	_ = v.BindEnv(FlagPlayers, "FARMFA_PLAYERS")
	_ = v.BindEnv(FlagThreshold, "FARMFA_THRESHOLD")
	_ = v.ReadInConfig()

	return c
}
