package router

import (
	"go-rest-api/controller"
	cv "go-rest-api/router/validator"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController) *echo.Echo {
	e := echo.New()

	e.Validator = &cv.CustomValidator{Validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	api.POST("/signup", uc.SignUp)
	api.POST("/login", uc.Login)
	api.POST("/logout", uc.Logout)

	return e
}
