package main

import (
	"go-pj-for-portfolio/controller"
	"go-pj-for-portfolio/db"
	"go-pj-for-portfolio/repository"
	"go-pj-for-portfolio/router"
	"go-pj-for-portfolio/usecase"
	"go-pj-for-portfolio/validator"
)

func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(*taskUsecase)

	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
}
