package server

import (
	"fmt"
	"github.com/SSSaaS/sssa-golang"
	"github.com/pquerna/otp/totp"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handle(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var s server
	s.sessions = map[string]*session{}

	e.POST("/shares", s.createShares)

	e.POST("/sessions", s.createSession)
	e.GET("/sessions/:id", s.readSession)
	e.POST("/sessions/:id/totp", s.getTOTPFromSession)
	e.POST("/sessions/:id/shares", s.joinShareWithSession)
}

type session struct {
	shares    []string
	private   string
	prefix    string
	threshold int

	mut sync.Mutex
}

type server struct {
	sessions map[string]*session
}

type CreateSharesRequest struct {
	SecretKey string   `json:"secret_key" form:"secret_key" query:"secret_key"`
	GpgKeys   []string `json:"gpg_keys" form:"gpg_keys" query:"gpg_keys"`

	Shares    int `json:"shares" form:"shares" query:"shares"`
	Threshold int `json:"threshold" form:"threshold" query:"threshold"`
}

type CreateSharesResponse struct {
	Shares []string `json:"shares"`
}

func (s *server) createShares(c echo.Context) error {
	var req CreateSharesRequest
	var resp CreateSharesResponse

	// Parse the request body into req
	if err := c.Bind(&req); err != nil {
		return err
	}

	// Split the secret key with SSS
	shares, err := sssa.Create(req.Threshold, req.Shares, req.SecretKey)
	if err != nil {
		return err
	}

	prefix, err := generateRandomString(5)
	if err != nil {
		return err
	}

	for i, share := range shares {
		resp.Shares = append(resp.Shares, fmt.Sprintf("%s$%d$%s", prefix, i, share))
	}

	return c.JSON(http.StatusOK, resp)
}

type CreateSessionRequest struct {
	Shares    int `json:"shares" form:"shares" query:"shares"`
	Threshold int `json:"threshold" form:"threshold" query:"threshold"`

	FirstShare string `json:"first_share" form:"first_share" query:"first_share"`
}

type CreateSessionResponse struct {
	ID      string `json:"id"`
	Private string `json:"private"`
}

func (s *server) createSession(c echo.Context) error {
	var req CreateSessionRequest
	var resp CreateSessionResponse

	// Parse the request body into req
	if err := c.Bind(&req); err != nil {
		return err
	}

	public, err := generateRandomString(25)
	if err != nil {
		return err
	}
	private, err := generateRandomString(25)
	if err != nil {
		return err
	}

	var newSession session

	firstShare, err := parseShare(req.FirstShare)
	if err != nil {
		return err
	}

	newSession.shares = append([]string{}, firstShare.share)

	newSession.prefix = firstShare.prefix
	newSession.private = private
	newSession.threshold = req.Threshold

	s.sessions[public] = &newSession

	resp.ID = public
	resp.Private = private

	return c.JSON(http.StatusOK, resp)
}

type ReadSessionResponse struct {
	Complete bool   `json:"complete"`
	Prefix   string `json:"prefix"`
}

func (s *server) readSession(c echo.Context) error {
	id := c.Param("id")

	var (
		resp ReadSessionResponse
		sess *session
		ok   bool
	)

	if sess, ok = s.sessions[id]; !ok {
		return c.String(http.StatusNotFound, "session id does not exist")
	}

	resp.Prefix = sess.prefix

	if len(sess.shares) >= sess.threshold {
		resp.Complete = true
	}

	return c.JSON(http.StatusOK, resp)
}

type GetTOTPRequest struct {
	Private string `json:"private" form:"private" query:"private"`
}

type GetTOTPResponse struct {
	TOTP string `json:"totp"`
}

func (s *server) getTOTPFromSession(c echo.Context) error {
	id := c.Param("id")

	var (
		req  GetTOTPRequest
		resp GetTOTPResponse
	)

	// Parse the request body into req
	if err := c.Bind(&req); err != nil {
		return err
	}

	var (
		sess *session
		ok   bool
	)

	if sess, ok = s.sessions[id]; !ok || sess.private != req.Private {
		return c.String(http.StatusBadRequest, "invalid ID or private key")
	}

	otpSecretKey, err := sssa.Combine(sess.shares)
	if err != nil {
		return c.String(http.StatusBadRequest, "some of the provided shares were invalid")
	}

	totp, err := totp.GenerateCode(otpSecretKey, time.Now())
	if err != nil {
		return c.String(http.StatusBadRequest, "could not generate TOTP")
	}

	resp.TOTP = totp

	return c.JSON(http.StatusOK, resp)
}

type JoinShareRequest struct {
	Share string `json:"share" form:"share" query:"share"`
}

func (s *server) joinShareWithSession(c echo.Context) error {
	id := c.Param("id")

	var req JoinShareRequest

	// Parse the request body into req
	if err := c.Bind(&req); err != nil {
		return err
	}

	share, err := parseShare(req.Share)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid share")
	}

	var (
		sess *session
		ok   bool
	)

	if sess, ok = s.sessions[id]; !ok {
		return c.String(http.StatusNotFound, "session id does not exist")
	}

	if sess.prefix != share.prefix {
		return c.String(http.StatusBadRequest, "provided prefix does not match current session")
	}

	if !sssa.IsValidShare(share.share) {
		return c.String(http.StatusBadRequest, "invalid share")
	}

	sess.mut.Lock()

	for _, existingShare := range sess.shares {
		if existingShare == share.share {
			sess.mut.Unlock()
			return c.String(http.StatusOK, "share already exists")
		}
	}

	sess.shares = append(sess.shares, share.share)
	sess.mut.Unlock()

	return c.String(http.StatusAccepted, "share joined")
}

type share struct {
	prefix string
	index  int
	share  string
}

func parseShare(completeShare string) (*share, error) {
	var parsed share

	parts := strings.Split(completeShare, "$")
	if n := len(parts); n != 3 {
		return nil, fmt.Errorf("split returned %d, wanted %d", n, 3)
	}

	parsed.prefix = parts[0]
	i, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid index: %w", err)
	}
	parsed.index = i
	parsed.share = parts[2]

	return &parsed, err
}
