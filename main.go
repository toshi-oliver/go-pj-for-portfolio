package main

import (
	"go-pj-for-portfolio/controller"
	"go-pj-for-portfolio/db"
	"go-pj-for-portfolio/repository"
	"go-pj-for-portfolio/router"
	"go-pj-for-portfolio/usecase"
)

func main() {
	db := db.NewDB()
	// userValidator := validator.NewUserValidator()
	// taskValidator := validator.NewTaskValidator()
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)
	e := router.NewRouter(userController)
	e.Logger.Fatal(e.Start(":8080"))
}
