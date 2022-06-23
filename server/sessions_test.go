package server_test

import (
	"github.com/borgoat/farmfa/server"
	"github.com/borgoat/farmfa/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	createSessionBody = `
{
	"toc_zero": {
		"group_id": "7GCUCI2Y",
		"group_size": 5,
		"group_threshold": 2,
		"share": "C2iCgb3pRfxPJw2a7od8p4ShkhrDWAm/Dt6ioQNAVFPZ",
		"toc_id": "5oaAUX9b6aBE"
	}
}
`
)

func TestCreateSession(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(createSessionBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := server.New(session.NewOracle(session.NewInMemoryStore()))

	if assert.NoError(t, h.CreateSession(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
