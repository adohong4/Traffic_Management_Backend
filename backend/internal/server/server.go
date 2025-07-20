package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

const (
	maxHeaderBytes = 1 << 20 // 1MB
	ctxTimeout     = 5       // time wait to shutdown (s)
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	db     *sqlx.DB
	logger logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, logger: logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * time.Duration(s.cfg.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.cfg.Server.WriteTimeout),
		MaxHeaderBytes: maxHeaderBytes,
	}

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	s.logger.Info("Shutting down server...")

	// create context with timeout to shutdown
	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	// Shutdown server
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.logger.Errorf("Error during server shutdown: %v", err)
		return err
	}

	s.logger.Info("Server exited properly")
	return nil
}
