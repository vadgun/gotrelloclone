package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/board-service/handlers"
	"github.com/vadgun/gotrelloclone/board-service/infra/config"
	"github.com/vadgun/gotrelloclone/board-service/infra/logger"
	"github.com/vadgun/gotrelloclone/board-service/infra/metrics"
	"github.com/vadgun/gotrelloclone/board-service/repositories"
	"github.com/vadgun/gotrelloclone/board-service/routes"
	"github.com/vadgun/gotrelloclone/board-service/services"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	metrics.InitMetrics()

	boardRepo := repositories.NewBoardRepository()
	boardService := services.NewBoardService(boardRepo)
	boardHandler := handlers.NewBoardHandler(boardService)

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupBoardRoutes(router, boardHandler)

	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))

	logger.Log.Info("🚀 board-service corriendo en http://board-service:8080")
	router.Run(":8080")
}
