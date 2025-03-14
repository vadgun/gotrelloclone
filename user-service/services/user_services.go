package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/vadgun/gotrelloclone/user-service/config"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"github.com/vadgun/gotrelloclone/user-service/repositories"
)

// UserService maneja la lógica de negocio de usuarios.
type UserService struct {
	repo *repositories.UserRepository
}

// NewUserService crea una nueva instancia del servicio.
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registra un usuario con contraseña encriptada.
func (s *UserService) RegisterUser(name, email, password, phone string) error {
	// Verificar si el usuario ya existe
	existingUser, _ := s.repo.GetUserByEmail(email)
	if existingUser != nil {
		return errors.New("el usuario ya existe")
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Crear usuario
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Phone:    phone,
	}

	// Guardar usuario en la BD
	return s.repo.CreateUser(user)
}

// LoginUser autentica un usuario y genera un token JWT.
func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expira en 24h
	})

	// Firmar el token
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserByID busca un usuario por su ID en la base de datos.
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
