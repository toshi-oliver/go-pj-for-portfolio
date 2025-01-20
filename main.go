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
	taskRepository := repository.NewTaskRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController)
	e.Logger.Fatal(e.Start(":8080"))
}
