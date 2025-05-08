package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vadgun/gotrelloclone/user-service/handlers"
	"github.com/vadgun/gotrelloclone/user-service/models"
	servicemocks "github.com/vadgun/gotrelloclone/user-service/services/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestUserHandler_Register_Sucess(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)

	// Datos de entrada
	input := map[string]string{
		"name":     "JoseR",
		"email":    "success@example.com",
		"password": "123456",
		"phone":    "1234567890",
	}
	jsonValue, _ := json.Marshal(input)

	// Simulamos la respuesta del mock
	mockService.On("RegisterUser", input["name"], input["email"], input["password"], input["phone"], "member").
		Return("mocked-id", nil)

	// Preparar el servidor de prueba
	router := setupRouter(mockHandler)

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

func TestUserHandler_Register_InvalidInput(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)

	// Datos de entrada
	invalidInput := map[string]string{
		"name":     "JoseR",
		"email":    "invalid@example.com",
		"password": "12345",
		"phone":    "123456789",
	}
	jsonValue, _ := json.Marshal(invalidInput)

	mockService.On("RegisterUser", invalidInput["name"], invalidInput["email"], invalidInput["password"], invalidInput["phone"], "member").
		Return("", errors.New("Télefono o contraseña invalidos"))

	// Server
	router := setupRouter(mockHandler)

	// Ejecutar la solicitud
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar la respuesta invalida
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Télefono o contraseña invalidos")
}

func TestUserHandler_Register_EmailExists(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)

	// Server
	router := setupRouter(mockHandler)

	userd := map[string]string{
		"name":     "JoseR",
		"email":    "existing@example.com",
		"password": "123456",
		"phone":    "1234567890",
	}
	userdjsv, _ := json.Marshal(userd)
	expectedError := errors.New("Error al registrar usuario")

	mockService.On("RegisterUser", userd["name"], userd["email"], userd["password"], userd["phone"], "member").
		Return("", expectedError)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userdjsv))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar la respuesta invalida
	assert.Equal(t, http.StatusConflict, resp.Code)
	assert.Contains(t, resp.Body.String(), expectedError.Error())
	mockService.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)
	router := setupRouter(mockHandler)

	t.Run("login existoso", func(t *testing.T) {
		user := &models.User{Name: "Jose", Email: "jose@test.com", Role: "admin"}
		mockService.On("LoginUser", user.Email, "secretagent").Return("mocked-token", user, nil)
		body := map[string]string{
			"email":    "jose@test.com",
			"password": "secretagent",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "mocked-token")
	})

	t.Run("login invalido", func(t *testing.T) {
		mockService.On("LoginUser", "", "").Return("", nil, nil)
		body := map[string]string{
			"email":    "",
			"password": "",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "usuario y contraseña requeridos al loguearse")

	})

	t.Run("usuario no existe", func(t *testing.T) {
		user := &models.User{Email: "userexists@test.com", Password: "secretagent"}
		expectedError := errors.New("usuario no registrado")
		mockService.On("LoginUser", user.Email, user.Password).Return("", user, expectedError)
		body := map[string]string{
			"email":    "userexists@test.com",
			"password": "secretagent",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Contains(t, resp.Body.String(), "usuario no registrado")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestUserHandler_Profile(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)

	t.Run("usuario no autenticado", func(t *testing.T) {
		router := setupRouter(mockHandler)

		userID := "mocked-id"
		expectedError := errors.New("Usuario no autenticado")
		mockService.On("GetUserByID", userID).Return(nil, expectedError)
		body := map[string]string{
			"userID": "mocked-id",
		}
		jsonValue, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodGet, "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), expectedError.Error())
	})

	t.Run("usuario encontrado", func(t *testing.T) {
		router := setupRouter(mockHandler)

		user := &models.User{ID: primitive.NewObjectID(), Name: "Jose"}
		mockService.On("GetUserByID", user.ID.Hex()).Return(user, nil)

		router.GET("/users/profile", func(ctx *gin.Context) {
			ctx.Set("userID", user.ID.Hex())
			mockHandler.Profile(ctx)
		})

		req, _ := http.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Jose")
	})

	t.Run("usuario no encontrado", func(t *testing.T) {
		router := setupRouter(mockHandler)

		nonExistentID := primitive.NewObjectID().Hex()
		expectedError := errors.New("Usuario no encontrado")

		mockService.On("GetUserByID", nonExistentID).Return((*models.User)(nil), expectedError)

		router.GET("/users/profile", func(ctx *gin.Context) {
			ctx.Set("userID", nonExistentID)
			mockHandler.Profile(ctx)
		})

		req, _ := http.NewRequest(http.MethodGet, "/users/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "Usuario no encontrado")
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)
	router := setupRouter(mockHandler)

	t.Run("Usuario existente", func(t *testing.T) {
		user := &models.User{ID: primitive.NewObjectID(), Name: "Jose"}
		mockService.On("GetUserByID", user.ID.Hex()).Return(user, nil)

		req, _ := http.NewRequest("GET", "/users/"+user.ID.Hex(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Jose")
	})

	t.Run("Usuario no existe", func(t *testing.T) {

		mockService.On("GetUserByID", "notfoundID").
			Return((*models.User)(nil), errors.New("No se pudo obtener el usuario"))

		req, _ := http.NewRequest(http.MethodGet, "/users/notfoundID", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "No se pudo obtener el usuario")
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	mockService := new(servicemocks.UserServiceMock)
	logger := zap.NewNop()
	mockHandler := handlers.NewUserHandler(mockService, logger)
	router := setupRouter(mockHandler)

	t.Run("obtener usuarios", func(t *testing.T) {

		users := []models.User{
			{ID: primitive.NewObjectID(), Name: "Jose"},
			{ID: primitive.NewObjectID(), Name: "Roberto"},
		}
		mockService.On("GetAllUsers").Return(users, nil)

		req, _ := http.NewRequest("GET", "/users", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Jose")
		assert.Contains(t, resp.Body.String(), "Roberto")
	})

	t.Run("no se pudo obtener todos los usuarios", func(t *testing.T) {
		expectecError := errors.New("Error al obtener usuarios")
		mockService := new(servicemocks.UserServiceMock)
		mockHandler := handlers.NewUserHandler(mockService, logger)
		router := setupRouter(mockHandler)

		mockService.On("GetAllUsers").
			Return(([]models.User)(nil), expectecError)

		req, _ := http.NewRequest("GET", "/users", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), expectecError.Error())

	})

}
