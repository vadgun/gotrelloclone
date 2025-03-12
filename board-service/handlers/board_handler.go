package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/board-service/services"
)

type BoardHandler struct {
	service *services.BoardService
}

func NewBoardHandler(service *services.BoardService) *BoardHandler {
	return &BoardHandler{service}
}

func (h *BoardHandler) CreateBoard(ctx *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := ctx.Get("userID") // Obtenemos el ID del usuario autenticado

	board, err := h.service.CreateBoard(request.Name, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el tablero"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"board": board})
}

func (h *BoardHandler) GetBoards(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	boards, err := h.service.GetBoardsByUser(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los tableros"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"boards": boards})
}

func (h *BoardHandler) GetBoardByID(ctx *gin.Context) {
	boardID := ctx.Param("boardID")
	tasks, err := h.service.GetBoardByID(boardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el tablero"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
