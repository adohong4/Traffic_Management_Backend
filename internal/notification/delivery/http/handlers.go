package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/internal/notification"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type notificationHandlers struct {
	cfg            *config.Config
	notificationUC notification.UseCase
	logger         logger.Logger
}

func NewNotificationHandlers(cfg *config.Config, notificationUC notification.UseCase, log logger.Logger) notification.Handlers {
	return &notificationHandlers{cfg: cfg, notificationUC: notificationUC, logger: log}
}

func (h *notificationHandlers) CreateNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		notification := &models.Notification{}
		if err := c.Bind(notification); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := notification.PrepareCreate(); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := utils.GetUserFromCtx(c.Request().Context())
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		notification.CreatorId = user.Id

		if err := utils.ValidateStruct(c.Request().Context(), notification); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		result, err := h.notificationUC.CreateNotification(c.Request().Context(), notification)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h *notificationHandlers) UpdateNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.Notification{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		result, err := h.notificationUC.UpdateNotification(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h *notificationHandlers) DeleteNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.Notification{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		deletedNotification, err := h.notificationUC.DeleteNotification(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, deletedNotification)
	}
}

func (h *notificationHandlers) GetNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.notificationUC.GetNotification(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}

func (h *notificationHandlers) GetNotificationById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		notificationID, err := h.notificationUC.GetNotificationByID(ctx, UUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, notificationID)
	}
}

func (h *notificationHandlers) SearchNotificationByTitle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		title := c.QueryParam("title")
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		newList, err := h.notificationUC.SearchNotificationByTitle(ctx, title, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, newList)
	}
}
