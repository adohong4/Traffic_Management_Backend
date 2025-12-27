package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/internal/notification"
	"github.com/labstack/echo/v4"
)

func MapNotificationRoutes(notificationGroup *echo.Group, h notification.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	notificationGroup.POST("/create", h.CreateNotification(), mw.AuthJWTMiddleware(authUC, cfg))
	notificationGroup.PUT("/:id", h.UpdateNotification(), mw.AuthJWTMiddleware(authUC, cfg))
	notificationGroup.DELETE("/:id", h.DeleteNotification(), mw.AuthJWTMiddleware(authUC, cfg))
	notificationGroup.GET("/:id", h.GetNotificationById())
	notificationGroup.GET("/getAll", h.GetNotification())
	notificationGroup.GET("/search", h.SearchNotificationByTitle())

	// === USER-SPECIFIC ROUTES ===
	notificationGroup.GET("/me", h.GetMyNotifications(), mw.AuthJWTMiddleware(authUC, cfg))
	notificationGroup.GET("/me/:id", h.GetMyNotificationByID(), mw.AuthJWTMiddleware(authUC, cfg))
}
