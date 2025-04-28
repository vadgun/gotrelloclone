package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/vadgun/gotrelloclone/user-service/infra/config"
	"github.com/vadgun/gotrelloclone/user-service/infra/kafka"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
)

// UserService maneja la lógica de negocio de usuarios.
type UserService struct {
	repo   repositories.UserRepositoryInterface
	Logger *zap.Logger
	Kafka  *kafka.KafkaProducer
}

// NewUserService crea una nueva instancia del servicio.
func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger, kafkaProducer *kafka.KafkaProducer) *UserService {
	return &UserService{repo: repo, Logger: logger, Kafka: kafkaProducer}
}

// RegisterUser registra un usuario con contraseña encriptada.
func (s *UserService) RegisterUser(name, email, password, phone, role string) (string, error) {
	// Verificar si el usuario ya existe
	existingUser, _ := s.repo.GetUserByEmail(email)
	if existingUser != nil {
		return "", errors.New("el usuario ya existe")
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Crear usuario
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Phone:    phone,
		Role:     role,
	}

	// Guardar usuario en la BD
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return "", err
	}

	event := &models.UserCreatedEvent{
		ID: id,
	}

	err = s.Kafka.Publish(context.TODO(), "new-user", event)
	if err != nil {
		s.Logger.Warn("⚠️ Usuario creado pero no se pudo publicar a Kafka", zap.String("id", id))
	}

	s.Logger.Info("✅ Usuario registrado correctamente")

	return id, nil
}

// LoginUser autentica un usuario y genera un token JWT.
func (s *UserService) LoginUser(email, password string) (string, *models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", nil, err
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("usuario o contraseña incorrectos")
	}

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expira en 24h
		"role":    user.Role,
	})

	// Firmar el token
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", user, err
	}

	return tokenString, user, nil
}

// GetUserByID busca un usuario por su ID en la base de datos.
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// GetAllUsers
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}
