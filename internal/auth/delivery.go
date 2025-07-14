package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	Delete() echo.HandlerFunc
	FindByIdentityNO() echo.HandlerFunc
	GetUsers() echo.HandlerFunc
	GetMe() echo.HandlerFunc
}
