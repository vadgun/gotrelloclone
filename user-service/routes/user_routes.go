package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/middlewares"
)

// SetupUserRoutes configura las rutas del usuario.
func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		userRoutes.POST("/login", userHandler.Login)

		// 🔐 Rutas protegidas
		userRoutes.GET("/profile", middlewares.AuthMiddleware(), userHandler.Profile)
		userRoutes.GET("/:userID", middlewares.AuthMiddleware(), userHandler.GetUserByID)
	}

}
