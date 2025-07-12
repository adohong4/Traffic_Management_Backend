package server

import (
	authHttp "github.com/adohong4/driving-license/internal/auth/delivery/http"
	authRepository "github.com/adohong4/driving-license/internal/auth/repository"
	authUseCase "github.com/adohong4/driving-license/internal/auth/usecase"
	apiMiddlewares "github.com/adohong4/driving-license/internal/middleware"
	"github.com/labstack/echo/v4"
)

// Map Server Handler
func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init Repositories
	aRepo := authRepository.NewAuthRepository(s.db)

	// Init Usecase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)

	// Init Handler
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

}
