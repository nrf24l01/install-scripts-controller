package httpsrv

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleInstall(c echo.Context) error {
	q := c.QueryParams()
	if !s.keys.Valid(q.Get("key")) {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	id, err := strconv.ParseInt(q.Get("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid script id")
	}

	var script string
	err = s.db.QueryRow(`SELECT script FROM scripts WHERE id = ?`, id).Scan(&script)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "script not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read script")
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.Blob(http.StatusOK, "text/plain; charset=utf-8", []byte(script))
}
