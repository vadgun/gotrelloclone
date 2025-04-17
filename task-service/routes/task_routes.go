package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/task-service/handlers"
	"github.com/vadgun/gotrelloclone/task-service/middlewares"
)

func SetupTaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middlewares.AuthMiddleware())

	taskGroup.POST("", taskHandler.CreateTask)                      // Crear tarea - Kafka producer implementado
	taskGroup.GET("/board/:boardID", taskHandler.GetTasksByBoardID) // Obtener todas las tareas de un board - Implementar producer
	taskGroup.GET("/:taskID", taskHandler.GetTaskByID)              // Obtener una tarea específica - Implementar producer
	taskGroup.PUT("/:taskID", taskHandler.UpdateTask)               // Actualizar tarea - Implementar producer
	taskGroup.PUT("/:taskID/move", taskHandler.MoveTask)            // Mueve una tarea entre boards - Implementar producer
	taskGroup.PUT("/:taskID/assign", taskHandler.AssignTask)        // Asigna la tarea a un usuario - Implementar producer
	taskGroup.PUT("/:taskID/status", taskHandler.UpdateTaskStatus)  // Cambia el estado de una tarea - Implementar producer
	taskGroup.DELETE("/:taskID", taskHandler.DeleteTask)            // Eliminar tarea - Implementar producer

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/tasks", middlewares.IsRoleAllowed("admin"), taskHandler.GetAllUsers)
	}

}
