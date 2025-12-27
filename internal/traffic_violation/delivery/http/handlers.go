package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	trafficviolation "github.com/adohong4/driving-license/internal/traffic_violation"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TrafficViolationHandlers struct {
	cfg                *config.Config
	TrafficViolationUC trafficviolation.UseCase
	logger             logger.Logger
}

func NewTrafficViolationHandlers(cfg *config.Config, TrafficViolationUC trafficviolation.UseCase, logger logger.Logger) trafficviolation.Handlers {
	return &TrafficViolationHandlers{cfg: cfg, TrafficViolationUC: TrafficViolationUC, logger: logger}
}

func (h *TrafficViolationHandlers) CreateTrafficViolation() echo.HandlerFunc {
	return func(c echo.Context) error {
		n := &models.TrafficViolation{}
		if err := c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		createdNew, err := h.TrafficViolationUC.CreateTrafficViolation(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdNew)
	}
}

func (h *TrafficViolationHandlers) UpdateTrafficViolation() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.TrafficViolation{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		updateData, err := h.TrafficViolationUC.UpdateTrafficViolation(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updateData)
	}
}

func (h *TrafficViolationHandlers) DeleteTrafficViolation() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.TrafficViolation{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		updateData, err := h.TrafficViolationUC.DeleteTrafficViolation(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updateData)
	}
}

func (h *TrafficViolationHandlers) GetTrafficViolationById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		geDetails, err := h.TrafficViolationUC.GetTrafficViolationById(ctx, UUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, geDetails)
	}
}

func (h *TrafficViolationHandlers) GetAllTrafficViolation() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.TrafficViolationUC.GetAllTrafficViolation(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

func (h *TrafficViolationHandlers) SearchTrafficViolation() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		query, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.TrafficViolationUC.SearchTrafficViolation(ctx, c.QueryParam("vehicle_no"), query)

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

func (h *TrafficViolationHandlers) GetTrafficViolationStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		stats, err := h.TrafficViolationUC.GetTrafficViolationStats(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, stats)
	}
}

func (h *TrafficViolationHandlers) GetTrafficViolationStatusStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		stats, err := h.TrafficViolationUC.GetTrafficViolationStatusStats(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, stats)
	}
}

// @Summary      Get all traffic violations of authenticated user
// @Description  Returns paginated list of traffic violations on vehicles owned by current user
// @Tags         traffic-violation
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        size  query     int  false  "Page size"
// @Success      200   {object}  models.TrafficViolationList
// @Failure      401   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic-violation/my-violations [get]
func (h *TrafficViolationHandlers) GetMyViolations() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.TrafficViolationUC.GetMyViolations(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// @Summary      Get violations for a specific vehicle owned by user
// @Description  Returns violations for a vehicle that belongs to the current user
// @Tags         traffic-violation
// @Produce      json
// @Param        vehicle_id  path  string  true  "Vehicle Registration ID"
// @Param        page        query int     false "Page number"
// @Param        size        query int     false "Page size"
// @Success      200         {object}  models.TrafficViolationList
// @Failure      400         {object}  httpErrors.RestError
// @Failure      401         {object}  httpErrors.RestError
// @Failure      404         {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic-violation/my-vehicles/{vehicle_id}/violations [get]
func (h *TrafficViolationHandlers) GetViolationsByMyVehicle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleID, err := uuid.Parse(c.Param("vehicle_id"))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.TrafficViolationUC.GetViolationsByMyVehicle(ctx, vehicleID, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// @Summary      Get detail of a traffic violation related to user's vehicle
// @Description  Retrieves a single traffic violation if it belongs to a vehicle owned by the current user
// @Tags         traffic-violation
// @Produce      json
// @Param        id    path      string  true  "Traffic Violation ID (UUID)"
// @Success      200   {object}  models.TrafficViolation
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      404   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic-violation/my-violations/{id} [get]
func (h *TrafficViolationHandlers) GetMyTrafficViolationByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		violationID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		violation, err := h.TrafficViolationUC.GetMyTrafficViolationByID(ctx, violationID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, violation)
	}
}

// @Summary      Get traffic violations related to user's driving license
// @Description  Returns violations on vehicles owned by the person whose driving license has the same wallet address
// @Tags         traffic-violation
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        size  query     int  false  "Page size"
// @Success      200   {object}  models.TrafficViolationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic-violation/my-license-violations [get]
func (h *TrafficViolationHandlers) GetViolationsByMyLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.TrafficViolationUC.GetViolationsByMyLicense(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}
