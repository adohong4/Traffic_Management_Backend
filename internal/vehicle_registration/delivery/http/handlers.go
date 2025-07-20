package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	vehicleRegistration "github.com/adohong4/driving-license/internal/vehicle_registration"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type vehicleRegHandlers struct {
	cfg          *config.Config
	vehicleRegUC vehicleRegistration.UseCase
	logger       logger.Logger
}

func NewVehicleReqHandlers(cfg *config.Config, vehicleRegUC vehicleRegistration.UseCase, logger logger.Logger) vehicleRegistration.Handlers {
	return &vehicleRegHandlers{cfg: cfg, vehicleRegUC: vehicleRegUC, logger: logger}
}

// Create godoc
// @Summary Create vehicle registration
// @Description Create vehicle_registration handler
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Success 201 {object} models.vehicle_registration
// @Router /vehicleReg/create [post]
func (h vehicleRegHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		n := &models.VehicleRegistration{}
		if err := c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		createdVehicleRegistration, err := h.vehicleRegUC.CreateVehicleDoc(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdVehicleRegistration)
	}
}

// Update godoc
// @Summary Update vehicle registration
// @Description Update vehicle_registration handler
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.vehicle_registration
// @Router /vehicleReg/{id} [put]
func (h vehicleRegHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleRegUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.VehicleRegistration{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.ID = vehicleRegUUID

		updatedVehicleReg, err := h.vehicleRegUC.UpdateVehicleDoc(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedVehicleReg)
	}
}

// Delete godoc
// @Summary Delete vehicle registration
// @Description Delete by id vehicle_registration handler
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.vehicle_registration
// @Router /vehicleReg/{id} [Delete]
func (h vehicleRegHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleRegUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.VehicleRegistration{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.ID = vehicleRegUUID

		deletedVehicleReg, err := h.vehicleRegUC.DeleteVehicleDoc(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, deletedVehicleReg)
	}
}

// GetByID godoc
// @Summary Get by vehicle registration ID
// @Description Get by vehicle registration handler
// @Tags Nvehicle registration
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.vehicle_registration
// @Router /vehicleReg/{id} [get]
func (h vehicleRegHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleRegUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		getVehicleRegID, err := h.vehicleRegUC.GetVehicleByID(ctx, vehicleRegUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, getVehicleRegID)
	}
}

// GetAllVehicleReg godoc
// @Summary Get all vehicle registration
// @Description Get all vehicle registration with pagination
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.vehicle registration list
// @Router /vehicleReg/getAll [get]
func (h vehicleRegHandlers) GetAllVehicleReg() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.vehicleRegUC.GetVehicleDocs(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

// SearchByVehiclePlateNO godoc
// @Summary Search by VehiclePlateNO
// @Description Search vehicle registration by VehiclePlateNO
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} models.vehicle_registration_list
// @Router /vehicleReg/search [get]
func (h vehicleRegHandlers) SearchByVehiclePlateNO() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.vehicleRegUC.FindByVehiclePlateNO(ctx, c.QueryParam("vehicle_no"), pq)

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}
