package repository

import (
	"fmt"
	"go-pj-for-portfolio/model"
	"math"

	"gorm.io/gorm"
)

type ITaskRepository interface {
	GetTasksByPage(userId uint, taskPage uint) (model.TaskResponsePaginated, error)
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	// UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

func (tr *taskRepository) GetTasksByPage(userId uint, taskPage uint) (model.TaskResponsePaginated, error) {
	const countPerPage = 10

	//ページングのために全件数を取得する
	var total int64
	if err := tr.db.Model(&model.Task{}).
		Where("user_id = ?", userId).
		Count(&total).Error; err != nil {
		return model.TaskResponsePaginated{}, err
	}

	// 総件数から最後のページ数（lastPage）を計算（切り上げ）
	lastPage := uint(math.Ceil(float64(total) / float64(countPerPage)))

	offset := (taskPage - 1) * countPerPage
	// ページ外のデータを取得しようとしていたらエラーにする
	if int64(offset) > total {
		return model.TaskResponsePaginated{}, fmt.Errorf("ページ外のデータを取得しようとしています")
	}
	// 指定したページのデータを降順で取得する
	var tasks []model.Task
	if err := tr.db.
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(countPerPage).
		Offset(int(offset)).
		Find(&tasks).Error; err != nil {
		return model.TaskResponsePaginated{}, err
	}

	// model.Task から TaskResponse へ変換
	resTasks := make([]model.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		resTasks = append(resTasks, model.TaskResponse{
			ID:         task.ID,
			MenuTitle:  task.MenuTitle,
			MenuDetail: task.MenuDetail,
			CreatedAt:  task.CreatedAt,
			UpdatedAt:  task.UpdatedAt,
		})
	}

	response := model.TaskResponsePaginated{
		Tasks:       resTasks,
		CurrentPage: taskPage,
		LastPage:    lastPage,
		TotalCount:  total,
	}

	return response, nil
}

func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
// 	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected < 1 {
// 		return fmt.Errorf("object does not exist")
// 	}
// 	return nil
// }

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
