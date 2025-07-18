package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapGovAgencyRoutes(GovAgencyGroup *echo.Group, h GovAgencyHandlers, mw *middleware.MiddlewareManager, cfg *config.Config, authUC auth.UseCase) {
	GovAgencyGroup.POST("/create", h.CreateGovAgency())
	GovAgencyGroup.PUT("/:id", h.UpdateGovAgency())
	GovAgencyGroup.DELETE("/:id", h.DeleteGovAgency())
	GovAgencyGroup.GET("/:id", h.GetByID())
	GovAgencyGroup.GET("/getAll", h.GetAllGovAgency())
	GovAgencyGroup.GET("/serach", h.SearchByName())
}
