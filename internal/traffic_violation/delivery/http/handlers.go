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

// @Summary      Create a new traffic violation record
// @Description  Creates a new traffic violation entry. Typically used by authorities/officers.
// @Tags         traffic-violation
// @Accept       json
// @Produce      json
// @Param        violation  body      models.TrafficViolation  true  "Traffic violation data"
// @Success      201        {object}  models.TrafficViolation
// @Failure      400        {object}  httpErrors.RestError  "Invalid request or validation error"
// @Failure      401        {object}  httpErrors.RestError  "Unauthorized"
// @Failure      500        {object}  httpErrors.RestError  "Internal server error"
// @Security     JWT
// @Router       /traffic/create [post]
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

// @Summary      Update an existing traffic violation
// @Description  Updates details of a traffic violation by ID (e.g., status, fine amount).
// @Tags         traffic-violation
// @Accept       json
// @Produce      json
// @Param        id         path      string                   true  "Traffic Violation ID (UUID)"
// @Param        violation  body      models.TrafficViolation  true  "Updated violation data"
// @Success      200        {object}  models.TrafficViolation
// @Failure      400        {object}  httpErrors.RestError
// @Failure      401        {object}  httpErrors.RestError
// @Failure      404        {object}  httpErrors.RestError  "Violation not found"
// @Failure      500        {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/{id} [put]
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

// @Summary      Soft delete a traffic violation
// @Description  Marks a traffic violation as inactive (soft delete). Used for correction or cancellation.
// @Tags         traffic-violation
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Traffic Violation ID (UUID)"
// @Success      200  {object}  models.TrafficViolation
// @Failure      400  {object}  httpErrors.RestError
// @Failure      401  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/{id} [delete]
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

// @Summary      Get traffic violation details by ID
// @Description  Retrieves a single traffic violation record (admin/officer view).
// @Tags         traffic-violation
// @Produce      json
// @Param        id   path      string  true  "Traffic Violation ID (UUID)"
// @Success      200  {object}  models.TrafficViolation
// @Failure      400  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError  "Violation not found"
// @Failure      500  {object}  httpErrors.RestError
// @Router       /traffic/{id} [get]
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

// @Summary      Get all traffic violations (paginated)
// @Description  Returns a paginated list of all active traffic violations (admin view).
// @Tags         traffic-violation
// @Produce      json
// @Param        page  query     int  false  "Page number (default: 1)"
// @Param        size  query     int  false  "Page size (default: 10)"
// @Success      200   {object}  models.TrafficViolationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Router       /traffic/getAll [get]
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

// @Summary      Search traffic violations by vehicle plate number
// @Description  Searches for violations containing the given plate number (partial, case-insensitive).
// @Tags         traffic-violation
// @Produce      json
// @Param        vehicle_no  query     string  true   "Vehicle plate number (partial)"
// @Param        page        query     int     false  "Page number"
// @Param        size        query     int     false  "Page size"
// @Success      200         {object}  models.TrafficViolationList
// @Failure      400         {object}  httpErrors.RestError  "Missing vehicle_no"
// @Failure      500         {object}  httpErrors.RestError
// @Router       /traffic/search [get]
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

// @Summary      Get overall traffic violation statistics
// @Description  Returns total violations, total fine amount, paid and unpaid amounts.
// @Tags         traffic-violation
// @Produce      json
// @Success      200  {object}  models.TrafficViolationStats
// @Failure      500  {object}  httpErrors.RestError
// @Router       /traffic/stats [get]
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

// @Summary      Get traffic violation statistics by status
// @Description  Returns breakdown by status including overdue/not overdue counts and amounts.
// @Tags         traffic-violation
// @Produce      json
// @Success      200  {array}   models.TrafficViolationStatusStats
// @Failure      500  {object}   httpErrors.RestError
// @Router       /traffic-violation/stats/status [get]
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
// @Tags         User
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        size  query     int  false  "Page size"
// @Success      200   {object}  models.TrafficViolationList
// @Failure      401   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/me [get]
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
// @Tags         User
// @Produce      json
// @Param        vehicle_id  path  string  true  "Vehicle Registration ID"
// @Param        page        query int     false "Page number"
// @Param        size        query int     false "Page size"
// @Success      200         {object}  models.TrafficViolationList
// @Failure      400         {object}  httpErrors.RestError
// @Failure      401         {object}  httpErrors.RestError
// @Failure      404         {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/me/{vehicle_id}/vehicle [get]
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
// @Tags         User
// @Produce      json
// @Param        id    path      string  true  "Traffic Violation ID (UUID)"
// @Success      200   {object}  models.TrafficViolation
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      404   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/me/{id} [get]
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
// @Tags         User
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        size  query     int  false  "Page size"
// @Success      200   {object}  models.TrafficViolationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /traffic/violation [get]
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
