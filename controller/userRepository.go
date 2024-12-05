package controller

import (
	"github.com/labstack/echo/v4"
)

type UserRepository interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	RefreshToken(c echo.Context) error
}
