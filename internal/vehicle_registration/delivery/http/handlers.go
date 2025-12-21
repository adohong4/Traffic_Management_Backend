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
// @Success 201 {object} models.VehicleRegistration
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
// @Success 200 {object} models.VehicleRegistration
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

// @Summary Confirm blockchain storage
// @Description Update blockchain transaction hash and set on_blockchain to true
// @Tags vehicle registrationn
// @Accept json
// @Produce json
// @Param id path string true "Vehicle Registration ID"
// @Param request body http.ConfirmBlockchainRequest true "Blockchain confirmation details"
// @Success 200 {object} models.VehicleRegistration
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /vehicle/{id}/confirm-blockchain [put]
func (h *vehicleRegHandlers) ConfirmBlockchainStorage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		req := &models.ConfirmBlockchainRequest{}
		if err = c.Bind(req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		dl := &models.VehicleRegistration{
			ID:               UUID,
			BlockchainTxHash: req.BlockchainTxHash,
			OnBlockchain:     req.OnBlockchain,
		}

		updatedDL, err := h.vehicleRegUC.ConfirmBlockchainStorage(ctx, dl)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedDL)
	}
}

// Delete godoc
// @Summary Delete vehicle registration
// @Description Delete by id vehicle_registration handler
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.VehicleRegistration
// @Router /vehicle/{id} [Delete]
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
// @Success 200 {object} models.VehicleRegistration
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
// @Success 200 {object} models.VehicleRegistration
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
// @Success 200 {object} models.VehicleRegistration
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

// GetStatsByType godoc
// @Summary Get vehicle count by type
// @Description Get count by vehicle type
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Success 200 {object} models.VehicleTypeCounts
// @Router /vehicleReg/stats/type [get]
func (h vehicleRegHandlers) GetStatsByType() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		stats, err := h.vehicleRegUC.GetCountByType(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, stats)
	}
}

// GetStatsByBrand godoc
// @Summary Get top brands distribution
// @Description Get top 5 brands and others
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Success 200 {object} models.BrandCounts
// @Router /vehicleReg/stats/brand [get]
func (h vehicleRegHandlers) GetStatsByBrand() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		stats, err := h.vehicleRegUC.GetTopBrands(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, stats)
	}
}

// GetStatsByStatus godoc
// @Summary Thống kê trạng thái đăng kiểm của xe cơ giới
// @Description Statistics on the number of motor vehicles by status: valid, expired, and pending inspection (based on ExpiryDate and RegistrationDate). Excludes motorcycles, motorbikes, mopeds, bicycles, and electric motorbikes.
// @Tags vehicle registration
// @Accept json
// @Produce json
// @Success 200 {object} models.StatusCounts
// @Router /vehicleReg/stats/status [get]
func (h vehicleRegHandlers) GetStatsByStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		stats, err := h.vehicleRegUC.GetCountByStatus(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, stats)
	}
}
