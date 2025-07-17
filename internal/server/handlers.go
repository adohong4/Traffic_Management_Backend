package server

import (
	"net/http"

	vehicleRegHttp "github.com/adohong4/driving-license/internal/vehicle_registration/delivery/http"
	vehicleReqRepository "github.com/adohong4/driving-license/internal/vehicle_registration/repository"
	vehicleReqUseCase "github.com/adohong4/driving-license/internal/vehicle_registration/usecase"

	authHttp "github.com/adohong4/driving-license/internal/auth/delivery/http"
	authRepository "github.com/adohong4/driving-license/internal/auth/repository"
	authUseCase "github.com/adohong4/driving-license/internal/auth/usecase"

	apiMiddlewares "github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Map Server Handler
func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init Repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	vReRepo := vehicleReqRepository.NewVehicleDocRepository(s.db)

	// Init Usecase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)
	vReUC := vehicleReqUseCase.NewVehicleRegUseCase(s.cfg, vReRepo, s.logger)

	// Init Handler
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)
	vehiclerReqHandlers := vehicleRegHttp.NewVehicleReqHandlers(s.cfg, vReUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)

	// middleware
	e.Use(mw.RequestLoggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// CSRF middleware
	if s.cfg.Server.CSRF {
		//if false {
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "header:X-CSRF-Token",
			CookieName:     s.cfg.Cookie.Name,
			CookieMaxAge:   s.cfg.Cookie.MaxAge,
			CookieSecure:   s.cfg.Cookie.Secure,
			CookieHTTPOnly: s.cfg.Cookie.HTTPOnly,
		}))
	}

	// API v1
	v1 := e.Group("v1/api")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	vehicleReq := v1.Group("/vehicleReg")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw, s.cfg, authUC)
	vehicleRegHttp.MapVehicleRegistrationRoutes(vehicleReq, vehiclerReqHandlers, mw, s.cfg, authUC)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check request id: %s", utils.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
