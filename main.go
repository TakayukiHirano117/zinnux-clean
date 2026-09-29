package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/infra/repository"
	"go-rest-api/infra/shared"
	"go-rest-api/router"
	"go-rest-api/usecase"
)

func main() {
	db := db.NewDB()

	userRepository := repository.NewUserRepository(db)
	passwordGenerator := shared.NewPasswordGenerator()
	userUsecase := usecase.NewUserUsecase(userRepository, passwordGenerator)
	userController := controller.NewUserController(userUsecase)
	e := router.NewRouter(userController)
	e.Logger.Fatal(e.Start(":8080"))
}
