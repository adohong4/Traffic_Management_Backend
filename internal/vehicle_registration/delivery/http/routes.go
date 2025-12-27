package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	vehicleRegistration "github.com/adohong4/driving-license/internal/vehicle_registration"
	"github.com/labstack/echo/v4"
)

func MapVehicleRegistrationRoutes(vehicleRegGroup *echo.Group, h vehicleRegistration.Handlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	vehicleRegGroup.POST("/create", h.Create(), mw.AuthJWTMiddleware(authUC, cfg))
	vehicleRegGroup.PUT("/:id", h.Update(), mw.AuthJWTMiddleware(authUC, cfg))
	vehicleRegGroup.PUT("/:id/confirm-blockchain", h.ConfirmBlockchainStorage(), mw.AuthJWTMiddleware(authUC, cfg))
	vehicleRegGroup.DELETE("/:id", h.Delete(), mw.AuthJWTMiddleware(authUC, cfg))
	vehicleRegGroup.GET("/:id", h.GetByID())
	vehicleRegGroup.GET("/getAll", h.GetAllVehicleReg())
	vehicleRegGroup.GET("/search", h.SearchByVehiclePlateNO())
	vehicleRegGroup.GET("/stats/type", h.GetStatsByType())
	vehicleRegGroup.GET("/stats/brand", h.GetStatsByBrand())
	vehicleRegGroup.GET("/stats/status", h.GetStatsByStatus())

	// User
	vehicleRegGroup.GET("/me", h.GetMyVehicles(), mw.AuthJWTMiddleware(authUC, cfg))
	vehicleRegGroup.GET("/me/:id", h.GetMyVehicleByID(), mw.AuthJWTMiddleware(authUC, cfg))
}
