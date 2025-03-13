package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

func (s *TaskService) GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error) {
	return s.repo.GetTasksByBoardID(ctx, boardID)
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
func (s *TaskService) BoardExists(ctx context.Context, boardID, token string) (bool, error) {
	// 📌 Hacemos una llamada al BoardService para verificar si el BoardID existe
	url := fmt.Sprintf("http://board-service:8080/boards/%s", boardID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Realizar la solicitud con un cliente HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return resp.StatusCode == http.StatusOK, nil

}

func (s *TaskService) MoveTask(ctx context.Context, taskID, newBoardID string) error {
	return s.repo.UpdateTaskBoard(ctx, taskID, newBoardID)
}

func (s *TaskService) AssignTask(ctx context.Context, taskID, userID string) error {
	s.SendNotification(userID, "Te han asignado una nueva tarea")
	return s.repo.UpdateTaskAssignee(ctx, taskID, userID)
}

func (s *TaskService) UserExists(ctx context.Context, userID, token string) (bool, error) {
	// 📌 Hacemos una llamada al UserService para verificar si el usuario existe
	url := fmt.Sprintf("http://user-service:8080/users/%s", userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Realizar la solicitud con un cliente HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) error {
	task, _ := s.repo.GetTaskByID(ctx, taskID)
	s.SendNotification(task.AssigneeID, fmt.Sprintf("El estado de tu tarea '%s' ha cambiado a '%s'", task.Title, status))
	return s.repo.UpdateTaskStatus(ctx, taskID, status)
}

func (s *TaskService) SendNotification(userID, message string) {
	url := "http://notification-service:8080/notify"
	payload := map[string]string{"user_id": userID, "message": message}
	jsonData, _ := json.Marshal(payload)

	_, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	// Enviar a WebSocket
	wsURL := "ws://notification-service:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Println("Error al conectar a WebSocket:", err)
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error al enviar mensaje WebSocket:", err)
	}
}
