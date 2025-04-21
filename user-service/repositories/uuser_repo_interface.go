package repositories

import "github.com/vadgun/gotrelloclone/user-service/models"

type UserRepositoryInterface interface {
	CreateUser(user *models.User) (string, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}
