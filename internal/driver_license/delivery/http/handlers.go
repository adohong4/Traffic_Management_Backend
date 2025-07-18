package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	driverlicense "github.com/adohong4/driving-license/internal/driver_license"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DriverLicenseHandlers struct {
	cfg             *config.Config
	DriverLicenseUC driverlicense.UseCase
	logger          logger.Logger
}

func NewDriverLicenseHandlers(cfg *config.Config, DriverLicenseUC driverlicense.UseCase, logger logger.Logger) driverlicense.Handlers {
	return &DriverLicenseHandlers{cfg: cfg, DriverLicenseUC: DriverLicenseUC, logger: logger}
}

func (h *DriverLicenseHandlers) CreateDriverLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		n := &models.DrivingLicense{}
		if err := c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		createdNews, err := h.DriverLicenseUC.CreateDriverLicense(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdNews)
	}
}

func (h *DriverLicenseHandlers) UpdateDriverLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.DrivingLicense{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		updatedDriverLicense, err := h.DriverLicenseUC.UpdateDriverLicense(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedDriverLicense)
	}
}

func (h *DriverLicenseHandlers) DeleteDriverLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.DrivingLicense{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		deletedDriverLicense, err := h.DriverLicenseUC.DeleteDriverLicense(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, deletedDriverLicense)
	}
}

func (h *DriverLicenseHandlers) GetDriverLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.DriverLicenseUC.GetDriverLicense(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

func (h *DriverLicenseHandlers) GetDriverLicenseById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		driverlicenseUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		getDriverLicenseID, err := h.DriverLicenseUC.GetDriverLicenseById(ctx, driverlicenseUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, getDriverLicenseID)
	}
}

func (h *DriverLicenseHandlers) SearchByLicenseNo() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.DriverLicenseUC.SearchByLicenseNo(ctx, c.QueryParam("license_no"), pq)

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}
