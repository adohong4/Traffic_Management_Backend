package driverlicense

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateDriverLicense() echo.HandlerFunc
	UpdateDriverLicense() echo.HandlerFunc
	ConfirmBlockchainStorage() echo.HandlerFunc
	AddWalletAddress() echo.HandlerFunc
	DeleteDriverLicense() echo.HandlerFunc
	GetDriverLicense() echo.HandlerFunc
	GetDriverLicenseById() echo.HandlerFunc
	GetDriverLicenseByWalletAddress() echo.HandlerFunc
	SearchByLicenseNo() echo.HandlerFunc
	GetStatusDistribution() echo.HandlerFunc
	GetLicenseTypeDistribution() echo.HandlerFunc
	GetLicenseTypeStatusDistribution() echo.HandlerFunc
	GetCityStatusDistribution() echo.HandlerFunc

	GetMyDrivingLicense() echo.HandlerFunc
}
