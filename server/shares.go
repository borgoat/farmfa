package server

import (
	"fmt"
	"net/http"

	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/shares"
	"github.com/labstack/echo/v4"
)

func (s Server) CreateShares(ctx echo.Context) error {
	var req api.CreateSharesJSONRequestBody
	var resp api.TOTPShares

	// Parse the request body into req
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// TODO Validate received TOTP secret, if URL, take out Base32 secret

	// Split the secret key with SSS
	tokens, err := shares.Split(shares.TOTPSecret(req.TotpSecretKey), uint(req.Threshold), uint(req.Shares))
	if err != nil {
		// TODO Public HTTP error, log internal
		return fmt.Errorf("error while splitting secret into shares: %w", err)
	}

	resp.Shares = make([]api.UserShare, len(tokens))

	for i, token := range tokens {
		s, err := token.String()
		if err != nil {
			// TODO Public HTTP error, log internal
			return fmt.Errorf("error while converting token to string")
		}

		if req.EncryptionKeys != nil {
			// TODO Encrypt shares
		}

		resp.Shares[i] = api.UserShare{
			Share: s,
			// TODO Encrypt shares and set user accordingly
			//User:
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}
