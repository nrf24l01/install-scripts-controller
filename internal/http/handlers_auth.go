package httpsrv

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *Server) handleLogin(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if req.Password != s.cfg.Site.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid password"})
	}
	return c.JSON(http.StatusOK, loginResponse{Token: s.sessions.Create()})
}

func (s *Server) handleLogout(c echo.Context) error {
	tok := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")
	s.sessions.Delete(tok)
	return c.NoContent(http.StatusNoContent)
}
