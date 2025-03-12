package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/services"
)

// UserController maneja las peticiones HTTP de usuario.
type UserHandler struct {
	service *services.UserService
}

// NewUserController crea una nueva instancia del controlador.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register maneja el registro de usuarios.
func (c *UserHandler) Register(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone" binding:"required,min=10"`
	}

	// Validar la entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Registrar usuario
	err := c.service.RegisterUser(req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}

// Login maneja la autenticación de usuarios.
func (c *UserHandler) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validar la entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Autenticar usuario y generar token
	token, err := c.service.LoginUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Profile devuelve la información del usuario autenticado.
func (c *UserHandler) Profile(ctx *gin.Context) {
	// Obtener userID del contexto (seteado por el middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar usuario en la base de datos
	user, err := c.service.GetUserByID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UserHandler) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	fmt.Println(userID)
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
