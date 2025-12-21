package vehicleRegistration

import "github.com/labstack/echo/v4"

type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	ConfirmBlockchainStorage() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	GetAllVehicleReg() echo.HandlerFunc
	SearchByVehiclePlateNO() echo.HandlerFunc
	GetStatsByType() echo.HandlerFunc
	GetStatsByBrand() echo.HandlerFunc
	GetStatsByStatus() echo.HandlerFunc
}
