package httpsrv

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// spa serves the built frontend and falls back to index.html for unknown
// routes so client-side routing keeps working.
func (s *Server) spa(c echo.Context) error {
	p := strings.TrimPrefix(c.Request().URL.Path, "/")
	target := filepath.Join(s.webDist, filepath.FromSlash(p))

	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		return c.File(filepath.Join(s.webDist, "index.html"))
	}
	return c.File(target)
}
