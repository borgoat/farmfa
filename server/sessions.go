package server

import (
	"net/http"

	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateSession(ctx echo.Context) error {
	var (
		req  api.CreateSessionJSONRequestBody
		resp *api.SessionCredentials
		err  error
	)

	if err := ctx.Bind(&req); err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusBadRequest, api.DefaultError{})
	}

	resp, err = s.oracle.CreateSession(&req.TocZero)
	if err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusInternalServerError, api.DefaultError{})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetSession(ctx echo.Context, id string) error {
	var (
		resp *api.Session
		err  error
	)

	resp, err = s.oracle.GetSession(id)
	if err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusInternalServerError, api.DefaultError{})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostToc(ctx echo.Context, id string) error {
	var (
		req api.PostTocJSONRequestBody
		err error
	)

	if err := ctx.Bind(&req); err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusBadRequest, api.DefaultError{})
	}

	err = s.oracle.AddToc(id, req.EncryptedToc)
	if err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusInternalServerError, api.DefaultError{})
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) GenerateTotp(ctx echo.Context, id string) error {
	var (
		req  api.GenerateTotpJSONRequestBody
		resp api.TOTPCode
		err  error
	)

	if err = ctx.Bind(&req); err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusBadRequest, api.DefaultError{})
	}

	//sess, err := s.oracle.GetSession(id)
	//if err != nil {
	//	// TODO Better error handling
	//	// most likely the ID is invalid
	//	return ctx.JSON(http.StatusBadRequest, api.DefaultError{})
	//}
	//
	//if sess.Status == "pending" {
	//	return ctx.JSON(http.StatusOK, sess)
	//}

	kek := api.SessionKeyEncryptionKey(req)
	totp, err := s.oracle.GenerateTOTP(id, &kek)
	if err != nil {
		// TODO Better error handling
		return ctx.JSON(http.StatusInternalServerError, api.DefaultError{})
	}

	resp.Totp = totp

	return ctx.JSON(http.StatusOK, resp)
}
