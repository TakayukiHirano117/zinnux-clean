package router

import (
	"go-rest-api/controller"
	cv "go-rest-api/router/validator"

	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
)

func NewRouter(uc controller.IUserController) *echo.Echo {
	e := echo.New()

	e.Validator = &cv.CustomValidator{Validator: validator.New()}

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.Login)
	e.POST("/logout", uc.Logout)

	return e
}
