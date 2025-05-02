package services

import "github.com/vadgun/gotrelloclone/user-service/models"

type UserServiceInterface interface {
	RegisterUser(name, email, password, phone, role string) (string, error)
	LoginUser(email, password string) (string, *models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}
