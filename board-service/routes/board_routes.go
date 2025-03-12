package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/board-service/handlers"
	"github.com/vadgun/gotrelloclone/board-service/middlewares"
)

func SetupBoardRoutes(router *gin.Engine, handler *handlers.BoardHandler) {
	boardGroup := router.Group("/boards")
	boardGroup.Use(middlewares.AuthMiddleware())

	boardGroup.POST("", handler.CreateBoard)
	boardGroup.GET("", handler.GetBoards)
	boardGroup.GET("/:boardID", handler.GetBoardByID)
}
