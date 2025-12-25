package http

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/user"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/labstack/echo/v4"
)

type userHandlers struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewNotificationHandlers(cfg *config.Config, userUC user.UseCase, log logger.Logger) user.Handlers {
	return &userHandlers{cfg: cfg, userUC: userUC, logger: log}
}

func (h *userHandlers) CreateNotification() echo.HandlerFunc {
	return func(c echo.Context) error {

	}
}
