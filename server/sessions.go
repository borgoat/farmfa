package server

import (
	"fmt"
	"github.com/giorgioazzinnaro/farmfa/ptr"
	"github.com/giorgioazzinnaro/farmfa/sessions"
	"github.com/giorgioazzinnaro/farmfa/shares"
	"net/http"

	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/labstack/echo/v4"
)

func (s Server) CreateSession(ctx echo.Context) error {
	var (
		req  api.CreateSessionJSONRequestBody
		resp api.PrivateSession
	)

	// Parse the request body into req
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	firstToken, err := shares.Parse(req.FirstShare)
	if err != nil {
		// TODO HTTP error
		return fmt.Errorf("the provide first share is invalid: %w", err)
	}

	session, err := s.sm.Start(firstToken)
	if err != nil {
		// TODO HTTP error
		return err
	}

	resp = *session
	return ctx.JSON(http.StatusOK, resp)
}

func (s Server) GetSession(ctx echo.Context, id string) error {
	var resp api.Session

	session, err := s.sm.Status(sessions.SessionIdentifier(id))
	if err != nil {
		return ctx.String(http.StatusNotFound, "session id does not exist")
	}

	resp = *session
	return ctx.JSON(http.StatusOK, resp)
}

func (s Server) PostShare(ctx echo.Context, id string) error {
	var (
		req  api.PostShareJSONRequestBody
		resp api.Session
	)

	// Parse the request body into req
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	token, err := shares.Parse(req.Share)
	if err != nil {
		// TODO HTTP error handling
		return fmt.Errorf("provided token is invalid: %w", err)
	}

	err = s.sm.AddShare(sessions.SessionIdentifier(id), token)
	if err != nil {
		// TODO HTTP error handling
		return fmt.Errorf("failed to add share: %w", err)
	}

	session, err := s.sm.Status(sessions.SessionIdentifier(id))
	if err != nil {
		return ctx.String(http.StatusNotFound, "session id does not exist")
	}

	resp = *session
	return ctx.JSON(http.StatusOK, resp)
}

func (s Server) GenerateTotp(ctx echo.Context, id string) error {
	var (
		req  api.GenerateTotpJSONRequestBody
		resp api.TOTPCode
	)

	// Parse the request body into req
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	totp, err := s.sm.GenerateTOTP(sessions.SessionIdentifier(id))
	if err != nil {
		// TODO HTTP error
		return fmt.Errorf("failed to generate TOTP: %w", err)
	}

	resp.Totp = ptr.String(string(totp))

	return ctx.JSON(http.StatusOK, resp)
}
