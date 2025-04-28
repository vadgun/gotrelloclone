package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vadgun/gotrelloclone/user-service/infra/kafka"
	"github.com/vadgun/gotrelloclone/user-service/infra/logger"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"github.com/vadgun/gotrelloclone/user-service/services"
	"golang.org/x/crypto/bcrypt"
)

// Mock del repositorio que prueba la interacción a través de la abstracción correcta.
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)

	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepo) GetUserByID(userID string) (*models.User, error) {
	args := m.Called(userID)

	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}

	return nil, args.Error(1)

}

func (m *MockUserRepo) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	if user := args.Get(0); user != nil {
		return user.([]models.User), args.Error(1)
	}

	return nil, args.Error(1)

}

func TestUserService_LoginUser_Success(t *testing.T) {
	// Repositorio Mockeado de user-repository
	mockRepo := new(MockUserRepo)

	// Obtener el Logger para nuestro servicio
	logger.InitLogger()
	log := logger.Log

	// Obtener el productor de kafka para enviar el mensaje?
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "testing_user_service", log)

	// Servicio que recibe un repositorio de user-repository
	mockService := services.NewUserService(mockRepo, log, kafkaProducer)

	// Generar el hash para el usuario a logearse en el test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretagent"), bcrypt.DefaultCost)

	// Usuario con datos hardcodeados para mayor legibilidad del test
	user := &models.User{
		Name:     "Jose Roberto",
		Email:    "vadgun@yahoo.com",
		Password: string(hashedPassword),
		Role:     "member",
	}

	// Inicia una descripción de una expectativa del método especificado que se está llamando.
	mockRepo.On("GetUserByEmail", "vadgun@yahoo.com").Return(user, nil)

	// Obtiene la respuesta de nuestro servicio
	token, resultUser, err := mockService.LoginUser("vadgun@yahoo.com", "secretagent")

	// Comprobamos
	// Que no exista error
	assert.NoError(t, err)

	// Que el token no esta vacio
	assert.NotEmpty(t, token)

	// Que ambos emails sean iguales
	assert.Equal(t, user.Email, resultUser.Email)

	// Confirma que todo lo especificado con On y Return se llamó según lo previsto.
	// Las llamadas pueden haber ocurrido en cualquier orden
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 1)
}

func TestUserService_LoginUser_InvalidPassword(t *testing.T) {
	// Repositorio Mockeado de user-repository
	mockRepo := new(MockUserRepo)

	// Generar el hash para el usuario a logearse en el test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secretagent"), bcrypt.DefaultCost)

	// Compara el password generado con el del ususario
	err := bcrypt.CompareHashAndPassword([]byte("secretAgent"), hashedPassword)
	assert.EqualError(t, err, "crypto/bcrypt: hashedSecret too short to be a bcrypted password")

	mockRepo.AssertExpectations(t)

}

func TestUserService_LoginUser_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepo)
	// Obtener el Logger para nuestro servicio
	logger.InitLogger()
	log := logger.Log
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "testing_user_service", log)
	service := services.NewUserService(mockRepo, log, kafkaProducer)

	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

	token, resultUser, err := service.LoginUser("notfound@example.com", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, resultUser)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 1)
}

func TestUserService_RegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	// Obtener el Logger para nuestro servicio
	logger.InitLogger()
	log := logger.Log
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "testing_user_service", log)
	service := services.NewUserService(mockRepo, log, kafkaProducer)

	t.Run("usuario creado correctamente", func(t *testing.T) {
		email := "test@example.com"
		mockRepo.On("GetUserByEmail", email).Return((*models.User)(nil), errors.New("not found"))
		mockRepo.On("CreateUser", mock.Anything).Return("12345", nil)

		userID, err := service.RegisterUser("Test", email, "password", "123456789", "user")

		assert.NotEmpty(t, userID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("usuario ya existe", func(t *testing.T) {
		email := "existing@example.com"
		mockRepo.On("GetUserByEmail", email).Return(&models.User{Email: email}, nil)

		userID, err := service.RegisterUser("Test", email, "password", "123456789", "user")

		assert.Empty(t, userID)
		assert.EqualError(t, err, "el usuario ya existe")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error al crear usuario cuando la bd falle", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := services.NewUserService(mockRepo, log, kafkaProducer)

		email := "fail@example.com"
		mockRepo.On("GetUserByEmail", email).Return((*models.User)(nil), errors.New("not found"))
		mockRepo.On("CreateUser", mock.Anything).Return("", errors.New("db error"))

		id, err := service.RegisterUser("Test", email, "password", "123456789", "user")

		assert.Empty(t, id)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepo) // Obtener el Logger para nuestro servicio
	logger.InitLogger()
	log := logger.Log
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "testing_user_service", log)
	service := services.NewUserService(mockRepo, log, kafkaProducer)

	t.Run("usuario encontrado", func(t *testing.T) {
		userID := "12345"
		expectedUser := &models.User{Email: "user@example.com"}

		mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)

		user, err := service.GetUserByID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("usuario no encontrado", func(t *testing.T) {
		userID := "notfound"
		mockRepo.On("GetUserByID", userID).Return(nil, errors.New("not found"))

		user, err := service.GetUserByID(userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepo) // Obtener el Logger para nuestro servicio
	logger.InitLogger()
	log := logger.Log
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "testing_user_service", log)
	service := services.NewUserService(mockRepo, log, kafkaProducer)

	expectedUsers := []models.User{}

	mockRepo.On("GetAllUsers").Return(expectedUsers, nil)
	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.EqualValues(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "GetAllUsers", 1)

}
