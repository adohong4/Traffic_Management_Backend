package http

import (
	"database/sql"
	"errors"
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

// @Summary Create a new driving license
// @Description Create a new driving license entry
// @Tags DrivingLicense
// @Accept json
// @Produce json
// @Param driving_license body models.DrivingLicense true "Driving License object"
// @Success 201 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /licenses/create [post]
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

// @Summary Update a driving license
// @Description Update an existing driving license by ID
// @Tags DrivingLicense
// @Accept json
// @Produce json
// @Param id path string true "Driving License ID"
// @Param driving_license body models.DrivingLicense true "Driving License object"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /licenses/{id} [put]
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

// @Summary Confirm blockchain storage
// @Description Update blockchain transaction hash and set on_blockchain to true
// @Tags DrivingLicense
// @Accept json
// @Produce json
// @Param id path string true "Driving License ID"
// @Param request body http.ConfirmBlockchainRequest true "Blockchain confirmation details"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /driver-license/{id}/confirm-blockchain [put]
func (h *DriverLicenseHandlers) ConfirmBlockchainStorage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		req := &ConfirmBlockchainRequest{}
		if err = c.Bind(req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		dl := &models.DrivingLicense{
			Id:               UUID,
			BlockchainTxHash: req.BlockchainTxHash,
			OnBlockchain:     req.OnBlockchain,
		}

		updatedDL, err := h.DriverLicenseUC.ConfirmBlockchainStorage(ctx, dl)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedDL)
	}
}

// @Summary Add wallet address
// @Description Add or update wallet address for a driving license
// @Tags DrivingLicense
// @Accept json
// @Produce json
// @Param id path string true "Driving License ID"
// @Param request body http.AddWalletRequest true "Wallet address details"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /driver-license/{id}/add-wallet [put]
func (h *DriverLicenseHandlers) AddWalletAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		req := &AddWalletRequest{}
		if err = c.Bind(req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := utils.ValidateStruct(ctx, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		dl := &models.DrivingLicense{
			Id:            UUID,
			WalletAddress: req.WalletAddress,
		}

		updatedDL, err := h.DriverLicenseUC.AddWalletAddress(ctx, dl)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedDL)
	}
}

// @Summary Delete a driving license
// @Description Soft delete a driving license by ID
// @Tags DrivingLicense
// @Accept json
// @Produce json
// @Param id path string true "Driving License ID"
// @Param driving_license body models.DrivingLicense true "Driving License object (for modifier)"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /licenses/{id} [delete]
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

// @Summary Get all driving licenses
// @Description Get a paginated list of all active driving licenses
// @Tags DrivingLicense
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} models.DrivingLicenseList
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/getAll [get]
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

// @Summary Get driving license by ID
// @Description Get a single driving license by its ID
// @Tags DrivingLicense
// @Produce json
// @Param id path string true "Driving License ID"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 404 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/{id} [get]
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

// @Summary Get driving license by Wallet Address
// @Description Get a single driving license by its Wallet Address
// @Tags DrivingLicense
// @Produce json
// @Param address path string true "Driving License Wallet Address"
// @Success 200 {object} models.DrivingLicense
// @Failure 400 {object} httpErrors.RestError
// @Failure 404 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/{address} [get]
func (h *DriverLicenseHandlers) GetDriverLicenseByWalletAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		walletAddress := c.Param("address")
		if walletAddress == "" {
			err := echo.NewHTTPError(http.StatusBadRequest, "wallet address is required")
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		result, err := h.DriverLicenseUC.GetDriverLicenseByWalletAddress(ctx, walletAddress)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, result)
	}
}

// @Summary Search driving licenses by license number
// @Description Search for driving licenses by license number with pagination
// @Tags DrivingLicense
// @Produce json
// @Param license_no query string true "License number to search"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} models.DrivingLicenseList
// @Failure 400 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/search [get]
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

// @Summary Get status distribution
// @Description Get count of driving licenses by status (active, pending, expired, pause, revoke)
// @Tags DrivingLicense
// @Produce json
// @Success 200 {object} models.StatusDistributionResponse
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/stats/status [get]
func (h *DriverLicenseHandlers) GetStatusDistribution() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		dist, err := h.DriverLicenseUC.GetStatusDistribution(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, dist)
	}
}

// @Summary Get license type distribution
// @Description Get count of driving licenses by license type (A1, B1, B2, C, ...)
// @Tags DrivingLicense
// @Produce json
// @Success 200 {object} models.LicenseTypeDistributionResponse
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/stats/license-type [get]
func (h *DriverLicenseHandlers) GetLicenseTypeDistribution() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		dist, err := h.DriverLicenseUC.GetLicenseTypeDistribution(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, dist)
	}
}

// @Summary Get detailed distribution by license type and status
// @Description Get count of driving licenses grouped by license_type, with breakdown by status (active, expired, pause, pending, revoke)
// @Tags DrivingLicense
// @Produce json
// @Success 200 {object} models.LicenseTypeDetailDistributionResponse
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/stats/license-type-detail [get]
func (h *DriverLicenseHandlers) GetLicenseTypeStatusDistribution() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		dist, err := h.DriverLicenseUC.GetLicenseTypeStatusDistribution(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, dist)
	}
}

// @Summary Get detailed distribution by owner city and status
// @Description Get count of driving licenses grouped by owner_city (province/city), with breakdown by status
// @Tags DrivingLicense
// @Produce json
// @Success 200 {object} models.CityDetailDistributionResponse
// @Failure 500 {object} httpErrors.RestError
// @Router /licenses/stats/city-detail [get]
func (h *DriverLicenseHandlers) GetCityStatusDistribution() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		dist, err := h.DriverLicenseUC.GetCityStatusDistribution(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, dist)
	}
}

// Confirm Blockchain Request
type ConfirmBlockchainRequest struct {
	BlockchainTxHash string `json:"blockchain_txhash" validate:"required"`
	OnBlockchain     bool   `json:"on_blockchain"`
}

// Add Wallet Address into Record
type AddWalletRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
}

// @Summary Get my driving license
// @Description Get detailed information of the driving license associated with the authenticated user (via wallet address)
// @Tags User
// @Produce json
// @Success 200 {object} models.DrivingLicense
// @Failure 401 {object} httpErrors.RestError
// @Failure 404 {object} httpErrors.RestError
// @Failure 500 {object} httpErrors.RestError
// @Security JWT
// @Router /licenses/me [get]
func (h *DriverLicenseHandlers) GetMyDrivingLicense() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.NewUnauthorizedError(err)))
		}

		if user.IdentityNo == "" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Không tìm thấy wallet address liên kết với tài khoản",
			})
		}

		dl, err := h.DriverLicenseUC.GetDriverLicenseByLicenseNO(ctx, user.IdentityNo)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusNotFound, map[string]string{
					"message": "Không tìm thấy bằng lái xe liên kết với tài khoản của bạn",
				})
			}
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, dl)
	}
}
