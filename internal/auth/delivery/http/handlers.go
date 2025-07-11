package http

import (
	"net/http"

	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/labstack/echo/v4"
)

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
}

func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		IdentityNO string `json:"identity_no" db:"identity_no" validate:"required,lte=20"`
		Password   string `json:"password,omitempty" db:"password" validate:"omitempty,gte=6"`
	}
	return func(c echo.Context) error {
		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		userWithToken, err := h.authUC.Login(ctx, &models.User{
			IdentityNo: login.IdentityNO,
			Password:   login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}
