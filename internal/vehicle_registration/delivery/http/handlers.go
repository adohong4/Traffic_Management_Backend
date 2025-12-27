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
// @Summary      Create a new vehicle registration
// @Description  Creates a new vehicle registration record. The vehicle plate number must be unique.
// @Tags         vehicle-registration
// @Accept       json
// @Produce      json
// @Param        vehicle  body      models.VehicleRegistration  true  "Vehicle registration data"
// @Success      201      {object}  models.VehicleRegistration
// @Failure      400      {object}  httpErrors.RestError
// @Failure      401      {object}  httpErrors.RestError
// @Failure      500      {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/create [post]
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
// @Summary      Update an existing vehicle registration
// @Description  Updates vehicle registration details by ID. Only provided fields are updated.
// @Tags         vehicle-registration
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Vehicle Registration ID (UUID)"
// @Param        vehicle  body      models.VehicleRegistration  true  "Updated vehicle registration data"
// @Success      200      {object}  models.VehicleRegistration
// @Failure      400      {object}  httpErrors.RestError
// @Failure      401      {object}  httpErrors.RestError
// @Failure      404      {object}  httpErrors.RestError
// @Failure      500      {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/{id} [put]
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

// ConfirmBlockchainStorage godoc
// @Summary      Confirm blockchain storage for a vehicle registration
// @Description  Updates the blockchain transaction hash and sets on_blockchain to true after successful storage.
// @Tags         vehicle-registration
// @Accept       json
// @Produce      json
// @Param        id       path      string                              true  "Vehicle Registration ID (UUID)"
// @Param        request  body      models.ConfirmBlockchainRequest     true  "Blockchain confirmation details"
// @Success      200      {object}  models.VehicleRegistration
// @Failure      400      {object}  httpErrors.RestError
// @Failure      401      {object}  httpErrors.RestError
// @Failure      500      {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/{id}/confirm-blockchain [put]
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
// @Summary      Soft delete a vehicle registration
// @Description  Marks a vehicle registration as inactive (soft delete).
// @Tags         vehicle-registration
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Vehicle Registration ID (UUID)"
// @Success      200  {object}  models.VehicleRegistration
// @Failure      400  {object}  httpErrors.RestError
// @Failure      401  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/{id} [delete]
func (h vehicleRegHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleRegUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.VehicleRegistration{ID: vehicleRegUUID}

		deletedVehicleReg, err := h.vehicleRegUC.DeleteVehicleDoc(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, deletedVehicleReg)
	}
}

// GetByID godoc
// @Summary      Get vehicle registration by ID
// @Description  Retrieves a single active vehicle registration record by its UUID.
// @Tags         vehicle-registration
// @Produce      json
// @Param        id   path      string  true  "Vehicle Registration ID (UUID)"
// @Success      200  {object}  models.VehicleRegistration
// @Failure      400  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /vehicle/{id} [get]
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
// @Summary      List all vehicle registrations
// @Description  Returns a paginated list of active vehicle registrations.
// @Tags         vehicle-registration
// @Produce      json
// @Param        page   query     int  false  "Page number (default: 1)"
// @Param        size   query     int  false  "Page size (default: 10)"
// @Success      200    {object}  models.VehicleRegistrationList
// @Failure      400    {object}  httpErrors.RestError
// @Failure      500    {object}  httpErrors.RestError
// @Router       /vehicle/getAll [get]
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
// @Summary      Search vehicle registrations by plate number
// @Description  Searches for active vehicle registrations containing the given plate number (partial match, case-insensitive).
// @Tags         vehicle-registration
// @Produce      json
// @Param        vehicle_no  query     string  true   "Plate number (partial)"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        size        query     int     false  "Page size (default: 10)"
// @Success      200         {object}  models.VehicleRegistrationList
// @Failure      400         {object}  httpErrors.RestError
// @Failure      500         {object}  httpErrors.RestError
// @Router       /vehicle/search [get]
func (h vehicleRegHandlers) SearchByVehiclePlateNO() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		vehicleNo := c.QueryParam("vehicle_no")
		if vehicleNo == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "vehicle_no query parameter is required"})
		}

		newList, err := h.vehicleRegUC.FindByVehiclePlateNO(ctx, vehicleNo, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

// GetStatsByType godoc
// @Summary      Vehicle statistics by type
// @Description  Returns count of vehicles by type: xe đầu kéo, ô tô tải, ô tô con, xe khách, xe máy, and "khác" (others).
// @Tags         vehicle-registration
// @Produce      json
// @Success      200  {array}   models.CountItem
// @Failure      500  {object}  httpErrors.RestError
// @Router       /vehicle/stats/type [get]
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
// @Summary      Top vehicle brands distribution
// @Description  Returns the top 5 most common vehicle brands and groups the rest under "khác".
// @Tags         vehicle-registration
// @Produce      json
// @Success      200  {array}   models.CountItem
// @Failure      500  {object}  httpErrors.RestError
// @Router       /vehicle/stats/brand [get]
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
// @Summary      Registration inspection status statistics (motor vehicles only)
// @Description  Returns count of motor vehicles (excluding motorcycles, mopeds, bicycles, electric bikes) by inspection status:
// @Description  - Valid: ExpiryDate >= today
// @Description  - Expired: ExpiryDate < today
// @Description  - Pending: No registration date or expiry date yet (e.g., new vehicles)
// @Tags         vehicle-registration
// @Produce      json
// @Success      200  {array}   models.CountItem
// @Failure      500  {object}  httpErrors.RestError
// @Router       /vehicle/stats/status [get]
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

// @Summary      List all vehicles owned by authenticated user
// @Description  Returns paginated list of vehicle registrations where owner_id matches current user
// @Tags         vehicle-registration
// @Produce      json
// @Param        page  query     int  false  "Page number (default: 1)"
// @Param        size  query     int  false  "Page size (default: 10)"
// @Success      200   {object}  models.VehicleRegistrationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/my-vehicles [get]
func (h vehicleRegHandlers) GetMyVehicles() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.vehicleRegUC.GetMyVehicles(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// @Summary      Get detail of a vehicle owned by authenticated user
// @Description  Retrieves a single vehicle registration that belongs to the current user
// @Tags         vehicle-registration
// @Produce      json
// @Param        id    path      string  true  "Vehicle Registration ID (UUID)"
// @Success      200   {object}  models.VehicleRegistration
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      404   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /vehicle/my-vehicles/{id} [get]
func (h vehicleRegHandlers) GetMyVehicleByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		vehicleID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		vehicle, err := h.vehicleRegUC.GetMyVehicleByID(ctx, vehicleID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, vehicle)
	}
}
