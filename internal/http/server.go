package httpsrv

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"install-scripts-controller/internal/config"
)

type Server struct {
	cfg      *config.Config
	db       *sql.DB
	sessions *SessionStore
	keys     *KeyStore
	webDist  string
	e        *echo.Echo
}

func New(cfg *config.Config, db *sql.DB, webDist string) *Server {
	s := &Server{
		cfg:      cfg,
		db:       db,
		sessions: NewSessionStore(24 * time.Hour),
		keys:     NewKeyStore(cfg.InstallKeyTTL()),
		webDist:  webDist,
		e:        echo.New(),
	}
	s.e.HideBanner = true
	s.e.HidePort = true
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.Logger())
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	e := s.e
	e.POST("/api/login", s.handleLogin)
	e.POST("/api/logout", s.auth(s.handleLogout))
	e.GET("/api/scripts", s.auth(s.handleListScripts))
	e.GET("/api/scripts/:id", s.auth(s.handleGetScript))
	e.POST("/api/scripts", s.auth(s.handleCreateScript))
	e.DELETE("/api/scripts/:id", s.auth(s.handleDeleteScript))

	e.GET("/install", s.handleInstall)

	e.GET("/", s.spa)
	e.GET("/*", s.spa)
}

func (s *Server) Start() error {
	return s.e.Start(s.cfg.Server.Addr)
}

// Handler exposes the underlying http.Handler (used by tests and embedders).
func (s *Server) Handler() http.Handler {
	return s.e
}
