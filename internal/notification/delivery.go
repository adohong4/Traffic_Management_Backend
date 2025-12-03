package notification

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateNotification() echo.HandlerFunc
	UpdateNotification() echo.HandlerFunc
	DeleteNotification() echo.HandlerFunc
	GetNotification() echo.HandlerFunc
	GetNotificationById() echo.HandlerFunc
	SearchNotificationByTitle() echo.HandlerFunc
}
