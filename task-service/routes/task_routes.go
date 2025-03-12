package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/task-service/handlers"
	"github.com/vadgun/gotrelloclone/task-service/middlewares"
)

func SetupTaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middlewares.AuthMiddleware())

	taskGroup.POST("", taskHandler.CreateTask)                      // Crear tarea
	taskGroup.GET("/board/:boardID", taskHandler.GetTasksByBoardID) // Obtener todas las tareas de un board
	taskGroup.GET("/:taskID", taskHandler.GetTaskByID)              // Obtener una tarea específica
	taskGroup.PUT("/:taskID", taskHandler.UpdateTask)               // Actualizar tarea
	taskGroup.PUT("/:taskID/move", taskHandler.MoveTask)            // Mueve una tarea entre boards
	taskGroup.PUT("/:taskID/assign", taskHandler.AssignTask)        // Asigna la tarea a un usuario
	taskGroup.PUT("/:taskID/status", taskHandler.UpdateTaskStatus)  // Cambia el estado de una tarea
	taskGroup.DELETE("/:taskID", taskHandler.DeleteTask)            // Eliminar tarea

}
