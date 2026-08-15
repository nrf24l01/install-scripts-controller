package httpsrv

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleHealth(c echo.Context) error {
	if err := s.db.PingContext(c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "database unavailable")
	}
	return c.NoContent(http.StatusOK)
}
