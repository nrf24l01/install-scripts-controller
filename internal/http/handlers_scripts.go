package httpsrv

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Script struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Script      string `json:"script,omitempty"`
	InstallURL  string `json:"install_url"`
	CreatedAt   string `json:"created_at"`
}

type createScriptRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Script      string `json:"script"`
}

func (s *Server) handleListScripts(c echo.Context) error {
	rows, err := s.db.Query(`SELECT id, name, description, created_at FROM scripts ORDER BY id DESC`)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list scripts")
	}
	defer rows.Close()

	scripts := []Script{}
	for rows.Next() {
		var sc Script
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Description, &sc.CreatedAt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read scripts")
		}
		sc.InstallURL = s.installURL(c, sc.ID)
		scripts = append(scripts, sc)
	}
	return c.JSON(http.StatusOK, scripts)
}

func (s *Server) handleGetScript(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid script id")
	}

	sc := Script{ID: id}
	err = s.db.QueryRow(
		`SELECT name, description, script, created_at FROM scripts WHERE id = ?`,
		id,
	).Scan(&sc.Name, &sc.Description, &sc.Script, &sc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "script not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read script")
	}
	sc.InstallURL = s.installURL(c, id)
	return c.JSON(http.StatusOK, sc)
}

func (s *Server) handleCreateScript(c echo.Context) error {
	var req createScriptRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if strings.TrimSpace(req.Script) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "script is required")
	}

	res, err := s.db.Exec(
		`INSERT INTO scripts (name, description, script) VALUES (?, ?, ?)`,
		req.Name, req.Description, req.Script,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create script")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create script")
	}

	var createdAt string
	if err := s.db.QueryRow(`SELECT created_at FROM scripts WHERE id = ?`, id).Scan(&createdAt); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read script")
	}

	return c.JSON(http.StatusCreated, Script{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		InstallURL:  s.installURL(c, id),
		CreatedAt:   createdAt,
	})
}

func (s *Server) handleDeleteScript(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid script id")
	}
	if _, err := s.db.Exec(`DELETE FROM scripts WHERE id = ?`, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete script")
	}
	return c.NoContent(http.StatusNoContent)
}

// installURL builds the public URL used to download a script's raw content
// with the install auth key embedded in the query string.
func (s *Server) installURL(c echo.Context, id int64) string {
	base := strings.TrimRight(s.cfg.Site.PublicURL, "/")
	if base == "" {
		base = c.Scheme() + "://" + c.Request().Host
	}
	return base + "/install?id=" + strconv.FormatInt(id, 10) +
		"&key=" + url.QueryEscape(s.keys.Current())
}
