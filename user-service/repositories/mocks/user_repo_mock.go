package repomocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vadgun/gotrelloclone/user-service/models"
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
