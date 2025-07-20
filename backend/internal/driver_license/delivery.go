package driverlicense

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateDriverLicense() echo.HandlerFunc
	UpdateDriverLicense() echo.HandlerFunc
	DeleteDriverLicense() echo.HandlerFunc
	GetDriverLicense() echo.HandlerFunc
	GetDriverLicenseById() echo.HandlerFunc
	SearchByLicenseNo() echo.HandlerFunc
}
