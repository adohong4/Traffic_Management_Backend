package trafficviolation

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateTrafficViolation() echo.HandlerFunc
	UpdateTrafficViolation() echo.HandlerFunc
	DeleteTrafficViolation() echo.HandlerFunc
	GetTrafficViolationById() echo.HandlerFunc
	GetAllTrafficViolation() echo.HandlerFunc
	SearchTrafficViolation() echo.HandlerFunc
	GetTrafficViolationStats() echo.HandlerFunc
	GetTrafficViolationStatusStats() echo.HandlerFunc
}
