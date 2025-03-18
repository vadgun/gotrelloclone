package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/task-service/models"
	"github.com/vadgun/gotrelloclone/task-service/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// 1️⃣ Crear tarea
func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userID, _ := ctx.Get("userID") // Obtenemos el ID del usuario autenticado

	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Split(authHeader, " ")

	// 📌 Validar si el BoardID existe antes de crear la tarea
	boardExists, err := h.service.BoardExists(ctx, task.BoardID, tokenString[1])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el BoardID"})
		return
	}
	if !boardExists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "El Board no existe"})
		return
	}

	id, err := h.service.CreateTask(ctx, &task, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la tarea"})
		return
	}
	// Usar h.service para seguir usando la arquitectura limpia RETO
	// 1ro intentar llegar a la funcion SendNotification enviando
	// valores o mensajes como la tarea en si, la cual viajara como un
	// record que contendra un string de toda la cadena en formato JSON
	// es nuestro caso, esto nos ayudara a capturar el evento(Tarea)
	// y guardarlo en nuestra base de datos.
	//
	// Para hacer que nuestros datos viajen a nuestro productor declarado
	// en nuestro h.service podemos pasarla a SendNotification varios datos
	// Topic , event-key , mensaje(tarea) para que sea consumido por
	// el consumer que esta escuchando en nuestro notification-service.
	// Convertir la tarea a JSON
	task.ID = id.(primitive.ObjectID)
	taskJSON, _ := json.Marshal(task)

	err = h.service.SendNotification(userID.(string), string(taskJSON), "task-events", "new-task")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando evento a Kafka"})
		return
	}

	// // Convertir la tarea a JSON
	// taskJSON, _ := json.Marshal(task)

	// Publicar evento en Kafka
	// err = kafka.ProduceMessage("task-events", "new-task", string(taskJSON))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Tarea creada y notificacion enviada exitosamente", "task": &task})
}

// 2️⃣ Obtener todas las tareas de un board
func (h *TaskHandler) GetTasksByBoardID(ctx *gin.Context) {
	boardID := ctx.Param("boardID")
	tasks, err := h.service.GetTasksByBoardID(ctx, boardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las tareas"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// 3️⃣ Obtener una tarea específica
func (h *TaskHandler) GetTaskByID(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	task, err := h.service.GetTaskByID(ctx, taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// 4️⃣ Actualizar tarea
func (h *TaskHandler) UpdateTask(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	var updatedData bson.M
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.service.UpdateTask(ctx, taskID, updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la tarea"})
		return
	}

	updatedDataJSON, _ := json.Marshal(updatedData)
	err = h.service.SendNotification(taskID, string(updatedDataJSON), "task-events", "updated-task")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando evento de update-task a Kafka"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tarea actualizada y notificacion enviada a kafka con éxito"})
}

// 5️⃣ Eliminar tarea
func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("taskID")

	err := h.service.DeleteTask(ctx, taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la tarea"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tarea eliminada correctamente"})
}

// 6️⃣ Mover tarea
func (h *TaskHandler) MoveTask(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	var request struct {
		NewBoardID string `json:"new_board_id"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Split(authHeader, " ")

	// 📌 Validar si el nuevo `BoardID` existe antes de mover la tarea
	boardExists, err := h.service.BoardExists(ctx, request.NewBoardID, tokenString[1])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el nuevo BoardID"})
		return
	}
	if !boardExists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "El Board destino no existe"})
		return
	}

	err = h.service.MoveTask(ctx, taskID, request.NewBoardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo mover la tarea"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tarea movida exitosamente"})
}

// 7️⃣ Asignar tarea a usuario
func (h *TaskHandler) AssignTask(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	var request struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 📌 Validar si el usuario existe en mongo-task
	userExists, err := h.service.UserExists(ctx, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar el usuario"})
		return
	}
	if !userExists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "El usuario no existe"})
		return
	}

	err = h.service.AssignTask(ctx, taskID, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo asignar la tarea"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tarea asignada exitosamente"})
}

// 8️⃣ Actualizar el status de una tarea
func (h *TaskHandler) UpdateTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	var request struct {
		Status models.TaskStatus `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		return
	}

	// Validar si el estado es válido
	validStatuses := map[models.TaskStatus]bool{
		models.TODO:       true,
		models.INPROGRESS: true,
		models.DONE:       true,
	}

	if !validStatuses[request.Status] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Estado no permitido"})
		return
	}

	err := h.service.UpdateTaskStatus(ctx, taskID, request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado en la base de datos"})
		return
	}

	taskStatusJSON, _ := json.Marshal(request.Status)

	err = h.service.SendNotification(taskID, string(taskStatusJSON), "task-events", "update-task-status")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando evento de update-task-status a Kafka"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Estado de la tarea y notificacion enviada a kafka con éxito"})
}
