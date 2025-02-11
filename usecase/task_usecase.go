package usecase

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"go-pj-for-portfolio/model"
	"go-pj-for-portfolio/repository"
	"go-pj-for-portfolio/validator"
	"log"
	"os"
)

type TaskUsecase struct {
	tr repository.ITaskRepository
}

func NewTaskUsecase(tr repository.ITaskRepository) *TaskUsecase {
	return &TaskUsecase{tr}
}

func (tu *TaskUsecase) GetTasksByPage(userId uint, taskPage uint) (model.TaskResponsePaginated, error) {
	response, err := tu.tr.GetTasksByPage(userId, taskPage)
	if err != nil {
		return model.TaskResponsePaginated{}, err
	}
	return response, nil
}

func (tu *TaskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	task := model.Task{}
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	resTask := model.TaskResponse{
		ID:         task.ID,
		MenuTitle:  task.MenuTitle,
		MenuDetail: task.MenuDetail,
		CreatedAt:  task.CreatedAt,
		UpdatedAt:  task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *TaskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	if err := validator.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}

	// FIXME: 以下のOpenAIの処理はGateway層に移動する。
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// APIキーを取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return model.TaskResponse{}, fmt.Errorf("環境変数にAPIキーが見つかりません")
	}

	// OpenAIクライアントを初期化
	llm, err := openai.New(openai.WithModel("gpt-4o-mini"))
	if err != nil {
		return model.TaskResponse{}, fmt.Errorf("OpenAIクライアントの作成に失敗: %v", err)
	}

	// プロンプトを生成
	prompt := fmt.Sprintf("System: \"ユーザーが入力した料理のレシピを教えてください\"\nUser: %s", task.MenuTitle)

	// ChatGPTにリクエストを送信し、応答を取得
	task.MenuDetail, err = llm.Call(context.Background(), prompt, llms.WithTemperature(0.5))
	if err != nil {
		return model.TaskResponse{}, fmt.Errorf("OpenAIへのリクエストに失敗: %v", err)
	}

	if err := tu.tr.CreateTask(&task); err != nil {
		return model.TaskResponse{}, err
	}

	resTask := model.TaskResponse{
		ID:         task.ID,
		MenuTitle:  task.MenuTitle,
		MenuDetail: task.MenuDetail,
		CreatedAt:  task.CreatedAt,
		UpdatedAt:  task.UpdatedAt,
	}

	return resTask, nil
}

// func (tu *TaskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
// 	if err := validator.TaskValidate(task); err != nil {
// 		return model.TaskResponse{}, err
// 	}

// 	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
// 		return model.TaskResponse{}, err
// 	}
// 	resTask := model.TaskResponse{
// 		ID:         task.ID,
// 		MenuTitle:  task.MenuTitle,
// 		MenuDetail: task.MenuDetail,
// 		CreatedAt:  task.CreatedAt,
// 		UpdatedAt:  task.UpdatedAt,
// 	}
// 	return resTask, nil
// }

func (tu *TaskUsecase) DeleteTask(userId uint, taskId uint) error {
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}
