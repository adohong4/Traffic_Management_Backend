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

// @Summary      Create a new notification
// @Description  Creates a new system notification. Only accessible by authenticated admin/staff users.
// @Tags         notification
// @Accept       json
// @Produce      json
// @Param        notification  body      models.Notification  true  "Notification data"
// @Success      200           {object}  models.Notification
// @Failure      400           {object}  httpErrors.RestError  "Invalid request body or validation error"
// @Failure      401           {object}  httpErrors.RestError  "Unauthorized"
// @Failure      500           {object}  httpErrors.RestError  "Internal server error"
// @Security     JWT
// @Router       /noti/create [post]
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

// @Summary      Update an existing notification
// @Description  Updates fields of a notification by ID. Only admin/staff can modify.
// @Tags         notification
// @Accept       json
// @Produce      json
// @Param        id            path      string               true  "Notification ID (UUID)"
// @Param        notification  body      models.Notification  true  "Updated notification data"
// @Success      200           {object}  models.Notification
// @Failure      400           {object}  httpErrors.RestError
// @Failure      401           {object}  httpErrors.RestError
// @Failure      404           {object}  httpErrors.RestError  "Notification not found"
// @Failure      500           {object}  httpErrors.RestError
// @Security     JWT
// @Router       /noti/{id} [put]
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

// @Summary      Soft delete a notification
// @Description  Marks a notification as inactive (soft delete). Only admin/staff.
// @Tags         notification
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Notification ID (UUID)"
// @Success      200  {object}  models.Notification
// @Failure      400  {object}  httpErrors.RestError
// @Failure      401  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Security     JWT
// @Router       /noti/{id} [delete]
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

// @Summary      Get paginated list of all notifications (admin view)
// @Description  Returns all active notifications with pagination. Mainly for admin dashboard.
// @Tags         notification
// @Produce      json
// @Param        page  query     int  false  "Page number (default: 1)"
// @Param        size  query     int  false  "Page size (default: 10)"
// @Success      200   {object}  models.NotificationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Router       /noti/getAll [get]
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

// @Summary      Get notification details by ID (admin view)
// @Description  Retrieves a single notification by ID. Accessible to admin or if user has permission.
// @Tags         notification
// @Produce      json
// @Param        id   path      string  true  "Notification ID (UUID)"
// @Success      200  {object}  models.Notification
// @Failure      400  {object}  httpErrors.RestError
// @Failure      404  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /noti/{id} [get]
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

// @Summary      Search notifications by title (admin view)
// @Description  Search active notifications containing the given title (case-insensitive partial match).
// @Tags         notification
// @Produce      json
// @Param        title  query     string  true   "Title keyword to search"
// @Param        page   query     int     false  "Page number (default: 1)"
// @Param        size   query     int     false  "Page size (default: 10)"
// @Success      200    {object}  models.NotificationList
// @Failure      400    {object}  httpErrors.RestError
// @Failure      500    {object}  httpErrors.RestError
// @Router       /noti/search [get]
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

// @Summary      Get list of notifications for authenticated user
// @Description  Returns notifications where target = "all" or target = "personal" with matching target_user (CCCD). Only shows notifications created after user account creation.
// @Tags         User
// @Produce      json
// @Param        page  query     int     false  "Page number (default: 1)"
// @Param        size  query     int     false  "Page size (default: 10)"
// @Success      200   {object}  models.NotificationList
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /noti/me [get]
func (h *notificationHandlers) GetMyNotifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.notificationUC.GetMyNotifications(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// @Summary      Get detail of a notification for user
// @Description  If the notification is personal and belongs to the user, it will be marked as read.
// @Tags         User
// @Produce      json
// @Param        id    path      string  true   "Notification ID"
// @Success      200   {object}  models.Notification
// @Failure      400   {object}  httpErrors.RestError
// @Failure      401   {object}  httpErrors.RestError
// @Failure      403   {object}  httpErrors.RestError
// @Failure      404   {object}  httpErrors.RestError
// @Failure      500   {object}  httpErrors.RestError
// @Security     JWT
// @Router       /noti/me/{id} [get]
func (h *notificationHandlers) GetMyNotificationByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		notificationID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		detail, err := h.notificationUC.GetMyNotificationByID(ctx, notificationID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, detail)
	}
}
