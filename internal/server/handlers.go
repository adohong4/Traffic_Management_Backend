package server

import (
	"net/http"

	_ "github.com/adohong4/driving-license/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	authHttp "github.com/adohong4/driving-license/internal/auth/delivery/http"
	authRepository "github.com/adohong4/driving-license/internal/auth/repository"
	authUseCase "github.com/adohong4/driving-license/internal/auth/usecase"

	govAgencyHttp "github.com/adohong4/driving-license/internal/gov_agency/delivery/http"
	govAgencyRepo "github.com/adohong4/driving-license/internal/gov_agency/repository"
	govAgencyUC "github.com/adohong4/driving-license/internal/gov_agency/usecase"

	driverLicenseHttp "github.com/adohong4/driving-license/internal/driver_license/delivery/http"
	driverLicenseRepo "github.com/adohong4/driving-license/internal/driver_license/repository"
	driverLicenseUseCase "github.com/adohong4/driving-license/internal/driver_license/usecase"

	vehicleRegHttp "github.com/adohong4/driving-license/internal/vehicle_registration/delivery/http"
	vehicleReqRepository "github.com/adohong4/driving-license/internal/vehicle_registration/repository"
	vehicleReqUseCase "github.com/adohong4/driving-license/internal/vehicle_registration/usecase"

	trafficVioHttp "github.com/adohong4/driving-license/internal/traffic_violation/delivery/http"
	trafficVioRepository "github.com/adohong4/driving-license/internal/traffic_violation/repository"
	trafficVioUseCase "github.com/adohong4/driving-license/internal/traffic_violation/usecase"

	newsHttp "github.com/adohong4/driving-license/internal/news/delivery/http"
	newsRepository "github.com/adohong4/driving-license/internal/news/repository"
	newsUseCase "github.com/adohong4/driving-license/internal/news/usecase"

	apiMiddlewares "github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Map Server Handler
func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init Repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	gRepo := govAgencyRepo.NewGovAgencyRepo(s.db)
	dRepo := driverLicenseRepo.NewDriverLicenseRepo(s.db)
	vReRepo := vehicleReqRepository.NewVehicleDocRepository(s.db)
	tRepo := trafficVioRepository.NewTrafficViolationRepo(s.db)
	newsRepo := newsRepository.NewNewsRepo(s.db)

	// Init Usecase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)
	goAgenUC := govAgencyUC.NewGovAgencyUseCase(s.cfg, gRepo, s.logger)
	dlUC := driverLicenseUseCase.NewDriverLicenseUseCase(s.cfg, dRepo, s.logger)
	vReUC := vehicleReqUseCase.NewVehicleRegUseCase(s.cfg, vReRepo, s.logger)
	tUC := trafficVioUseCase.NewTrafficViolationUseCase(s.cfg, tRepo, s.logger)
	newsUC := newsUseCase.NewNewsUseCase(s.cfg, newsRepo, s.logger)

	// Init Handler
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)
	govAgencyHandlers := govAgencyHttp.NewGovAgencyHandlers(s.cfg, goAgenUC, s.logger)
	driverLicenseHandlers := driverLicenseHttp.NewDriverLicenseHandlers(s.cfg, dlUC, s.logger)
	vehiclerReqHandlers := vehicleRegHttp.NewVehicleReqHandlers(s.cfg, vReUC, s.logger)
	trafficVioHandlers := trafficVioHttp.NewTrafficViolationHandlers(s.cfg, tUC, s.logger)
	newsHandlers := newsHttp.NewsHandlers(s.cfg, newsUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)

	// middleware
	e.Use(mw.RequestLoggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	//Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
	goAgencyGroup := v1.Group("/agency")
	driverLicenseGroup := v1.Group("/licenses")
	vehicleReqGroup := v1.Group("/vehicleReg")
	trafficVioGroup := v1.Group("/traffic")
	newsGroup := v1.Group("/news")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw, s.cfg, authUC)
	govAgencyHttp.MapGovAgencyRoutes(goAgencyGroup, govAgencyHandlers)
	driverLicenseHttp.MapDriverLicenseRoutes(driverLicenseGroup, driverLicenseHandlers, mw, s.cfg, authUC)
	vehicleRegHttp.MapVehicleRegistrationRoutes(vehicleReqGroup, vehiclerReqHandlers, mw, s.cfg, authUC)
	trafficVioHttp.MapTrafficViolationRoutes(trafficVioGroup, trafficVioHandlers, mw, s.cfg, authUC)
	newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw, authUC, s.cfg)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check request id: %s", utils.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
