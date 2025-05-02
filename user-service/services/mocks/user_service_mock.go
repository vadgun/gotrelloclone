package servicemocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vadgun/gotrelloclone/user-service/models"
)

type UserServiceMock struct {
	mock.Mock
}

func (sm *UserServiceMock) RegisterUser(name, email, password, phone, role string) (string, error) {
	args := sm.Called(name, email, password, phone, role)
	return args.String(0), args.Error(1)
}

func (sm *UserServiceMock) LoginUser(email, password string) (string, *models.User, error) {
	args := sm.Called(email, password)
	return args.String(0), args.Get(1).(*models.User), args.Error(2)
}

func (sm *UserServiceMock) GetUserByID(userID string) (*models.User, error) {
	args := sm.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (sm *UserServiceMock) GetAllUsers() ([]models.User, error) {
	args := sm.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
