package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/config"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
	"github.com/vadgun/gotrelloclone/user-service/routes"
	"github.com/vadgun/gotrelloclone/user-service/services"
)

func main() {
	// Iniciar conexión a MongoDB
	config.InitConfig()

	// Inicializar repositorio y servicio
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Configurar el servicio en modo producción
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()

	// Configurar rutas
	routes.SetupUserRoutes(router, userHandler)

	// Iniciar servidor en el puerto 8080
	log.Println("🚀 user-service corriendo en http://user-service:8080")
	router.Run(":8080")
}
