package controller

import (
	"go-pj-for-portfolio/model"
	"go-pj-for-portfolio/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

type TaskController struct {
	tu usecase.TaskUsecase
}

func NewTaskController(tu usecase.TaskUsecase) *TaskController {
	return &TaskController{tu}
}

func (tc *TaskController) GetTasksByPage(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	page := 1
	var err error
	if c.QueryParam("page") != "" {
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "不正なページが指定されました")
		}
	}

	tasksRes, err := tc.tu.GetTasksByPage(uint(userId.(float64)), uint(page))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasksRes)
}

func (tc *TaskController) GetTaskById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	taskRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskRes)
}

func (tc *TaskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task.UserId = uint(userId.(float64))
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, taskRes)
}

// func (tc *TaskController) UpdateTask(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	userId := claims["user_id"]
// 	id := c.Param("taskId")
// 	taskId, _ := strconv.Atoi(id)

// 	task := model.Task{}
// 	if err := c.Bind(&task); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}
// 	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, taskRes)
// }

func (tc *TaskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
