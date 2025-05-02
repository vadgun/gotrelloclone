package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	servicemocks "github.com/vadgun/gotrelloclone/user-service/services/mocks"
	"go.uber.org/zap"
)

func setupRouter(handler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.GET("/profile", handler.Profile)
	r.GET("/users/:userID", handler.GetUserByID)
	r.GET("/users", handler.GetAllUsers)
	return r
}

func TestRegister_Sucess(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	handler := handlers.NewUserHandler(mockService, logger)

	// Datos de entrada
	input := map[string]string{
		"name":     "Juan Perez",
		"email":    "juan@example.com",
		"password": "123456",
		"phone":    "1234567890",
	}
	jsonValue, _ := json.Marshal(input)

	// Simulamos la respuesta del mock
	mockService.On("RegisterUser", input["name"], input["email"], input["password"], input["phone"], "member").
		Return("mocked-id", nil)

	// Preparar el servidor de prueba
	router := setupRouter(handler)

	// Ejecutar la solicitud
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Usuario registrado correctamente")

	mockService.AssertExpectations(t)
}
