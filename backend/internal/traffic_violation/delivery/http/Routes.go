package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	trafficviolation "github.com/adohong4/driving-license/internal/traffic_violation"
	"github.com/labstack/echo/v4"
)

func MapTrafficViolationRoutes(trafficViolationGroup *echo.Group, h trafficviolation.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	trafficViolationGroup.POST("/create", h.CreateTrafficViolation(), mw.AuthJWTMiddleware(authUC, cfg))
	trafficViolationGroup.PUT("/:id", h.UpdateTrafficViolation(), mw.AuthJWTMiddleware(authUC, cfg))
	trafficViolationGroup.DELETE("/:id", h.DeleteTrafficViolation(), mw.AuthJWTMiddleware(authUC, cfg))
	trafficViolationGroup.GET("/:id", h.GetTrafficViolationById())
	trafficViolationGroup.GET("/getAll", h.GetAllTrafficViolation())
	trafficViolationGroup.GET("/search", h.SearchTrafficViolation())
}
