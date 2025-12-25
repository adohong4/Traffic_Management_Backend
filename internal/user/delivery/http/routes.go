package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/internal/user"
	"github.com/labstack/echo/v4"
)

func MapUserRoutes(userGroup *echo.Group, h user.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	userGroup.POST("/create", h.CreateUser(), mw.AuthJWTMiddleware(authUC, cfg))
}
