package services

import (
	"context"

	"github.com/vadgun/gotrelloclone/task-service/kafka"
	"github.com/vadgun/gotrelloclone/task-service/models"
	"github.com/vadgun/gotrelloclone/task-service/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

type TaskService struct {
	repo *repositories.TaskRepository
}

func NewTaskService(repo *repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task, userID string) (interface{}, error) {
	return s.repo.CreateTask(ctx, task, userID)
}

func (s *TaskService) GetTasksByBoardID(ctx context.Context, boardID string, page, limit int64) ([]models.Task, int64, error) {
	return s.repo.GetTasksByBoardID(ctx, boardID, page, limit)
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID string) (*models.Task, error) {
	return s.repo.GetTaskByID(ctx, taskID)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID string, updatedData bson.M) error {
	return s.repo.UpdateTask(ctx, taskID, updatedData)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	return s.repo.DeleteTask(ctx, taskID)
}

func (s *TaskService) MoveTask(ctx context.Context, taskID, newBoardID string) error {
	return s.repo.UpdateTaskBoard(ctx, taskID, newBoardID)
}

func (s *TaskService) AssignTask(ctx context.Context, taskID, userID string) error {
	return s.repo.UpdateTaskAssignee(ctx, taskID, userID)
}

func (s *TaskService) UserExists(ctx context.Context, userID string) (bool, error) {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) (err error) {
	if err := s.repo.UpdateTaskStatus(ctx, taskID, status); err != nil {
		return err
	}
	return err
}

func (s *TaskService) BoardExists(ctx context.Context, boardID string) (bool, error) {
	_, err := s.repo.GetBoardByID(boardID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TaskService) SendNotification(userID, message, topic, key string) error {
	err := kafka.ProduceMessage(userID, message, topic, key)
	return err
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAllTasks()
}
