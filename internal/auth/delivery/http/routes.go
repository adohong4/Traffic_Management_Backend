package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	authGroup.POST("/create", h.CreateUser())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/connectWallet", h.ConnectWallet())
	authGroup.POST("/logout", h.Logout())
	authGroup.PUT("/update/:id", h.Update())
	authGroup.GET("/find/", h.FindByIdentityNO())
	authGroup.DELETE("/delete/:id", h.Delete(), mw.AuthJWTMiddleware(authUC, cfg))
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:id", h.GetUserByID())
	authGroup.GET("/me", h.GetMe(), mw.AuthJWTMiddleware(authUC, cfg))
}
