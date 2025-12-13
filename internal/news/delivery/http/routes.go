package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/internal/news"
	"github.com/labstack/echo/v4"
)

func MapNewsRoutes(newsGroup *echo.Group, h news.Handlers, mw *middleware.MiddlewareManager, authUC auth.UseCase, cfg *config.Config) {
	newsGroup.POST("", h.Create(), mw.AuthJWTMiddleware(authUC, cfg))
	newsGroup.PUT("/:id", h.Update(), mw.AuthJWTMiddleware(authUC, cfg))
	newsGroup.DELETE("/:id", h.Delete(), mw.AuthJWTMiddleware(authUC, cfg))
	newsGroup.GET("/:id", h.FindById())
	newsGroup.GET("", h.FindAll())
}
