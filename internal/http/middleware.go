package httpsrv

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tok := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")
		if tok == "" || !s.sessions.Valid(tok) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return next(c)
	}
}
