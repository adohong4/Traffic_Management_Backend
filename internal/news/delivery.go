package news

import "github.com/labstack/echo/v4"

type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	FindById() echo.HandlerFunc
	FindAll() echo.HandlerFunc
}
