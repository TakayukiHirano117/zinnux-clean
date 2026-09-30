// Package controller
package controller

import (
	"go-rest-api/infra/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"min=0,max=255,required,email"`
	Password string `json:"password" validate:"min=8,max=100,required"`
}

type SignUpResponse struct {
	ID string `json:"id"`
}


func (uc *userController) SignUp(c echo.Context) error {
	signUpRequest := SignUpRequest{}
	if err := c.Bind(&signUpRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&signUpRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	signUpRequestDTO := usecase.SignUpRequestDTO{
		Email:    signUpRequest.Email,
		Password: signUpRequest.Password,
	}

	userRes, err := uc.uu.SignUp(&signUpRequestDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	signUpResponse := SignUpResponse{
		ID: string(userRes.ID),
	}
	return c.JSON(http.StatusCreated, signUpResponse)
}

func (uc *userController) Login(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func (uc *userController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
