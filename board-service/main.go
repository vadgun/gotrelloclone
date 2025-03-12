package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/board-service/config"
	"github.com/vadgun/gotrelloclone/board-service/handlers"
	"github.com/vadgun/gotrelloclone/board-service/repositories"
	"github.com/vadgun/gotrelloclone/board-service/routes"
	"github.com/vadgun/gotrelloclone/board-service/services"
)

func main() {
	config.InitConfig()

	boardRepo := repositories.NewBoardRepository()
	boardService := services.NewBoardService(boardRepo)
	boardHandler := handlers.NewBoardHandler(boardService)

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	routes.SetupBoardRoutes(router, boardHandler)

	log.Println("🚀 board-service corriendo en http://board-service:8080")
	router.Run(":8080")
}
