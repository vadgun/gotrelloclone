package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadgun/gotrelloclone/user-service/infra/metrics"
	"github.com/vadgun/gotrelloclone/user-service/services"
	"go.uber.org/zap"
)

// UserController maneja las peticiones HTTP de usuario.
type UserHandler struct {
	service services.UserServiceInterface
	Logger  *zap.Logger
}

// NewUserController crea una nueva instancia del controlador.
func NewUserHandler(service services.UserServiceInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, Logger: logger}
}

// Register maneja el registro de usuarios.
func (c *UserHandler) Register(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone" binding:"required,min=10"`
		Role     string
	}

	// Validar la entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// c.Logger.Info("❌ Error en el body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errors.New(" Télefono o contraseña invalidos").Error()})
		return
	}

	// Registrar usuario usando un rol por defecto
	id, err := c.service.RegisterUser(req.Name, req.Email, req.Password, req.Phone, "member")
	if err != nil {
		c.Logger.Info("❌ Error al registrar usuario", zap.Error(err))
		ctx.JSON(http.StatusConflict, gin.H{"error": "Error al registrar usuario"})
		return
	}

	// Metrica que lleva el numero de usuarios registrados
	metrics.UsersCreated.Inc()

	// Crear log personalizado
	c.Logger.Info("Guardando usuario en la base de datos", zap.String("endpoint", ctx.Request.URL.Path), zap.String("method", "POST"), zap.String("_id", id))

	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente", id: id})
}

// Login maneja la autenticación de usuarios.
func (c *UserHandler) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validar la entrada
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Logger.Info("❌ Usuario y contraseña requeridos al loguearse", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Autenticar usuario y generar token
	token, user, err := c.service.LoginUser(req.Email, req.Password)
	if err != nil {
		metrics.LoginAttempts.WithLabelValues("fail").Inc()
		c.Logger.Info("❌ Datos incorrectos o usuario no registrado", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("usuario no registrado").Error()})
		return
	}
	metrics.LoginAttempts.WithLabelValues("success").Inc()

	// Crear log personalizado
	c.Logger.Info("Usuario loggeado", zap.String("endpoint", ctx.Request.URL.Path), zap.String("method", "POST"), zap.String("user_email", req.Email))

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user.Name, "role": user.Role})
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
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el usuario"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
