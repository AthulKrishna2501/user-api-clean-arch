package services_test

import (
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestSingup(t *testing.T) {
	mockRepo := new(MockUserRepository)
	UserService := services.NewUserService(mockRepo)

	input := &models.SignupInput{
		UserName:    "JohnDoe",
		Email:       "johndoe@gmail.com",
		Password:    "johndoe123",
		PhoneNumber: "1234567890",
	}

	mockRepo.On("FindUserByEmail", input.Email).Return(nil, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	err := UserService.SignUp(input)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

}

func TestSingupUserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	input := &models.SignupInput{
		UserName:    "JohnDoe",
		Email:       "johndoegmail.com",
		Password:    "johndoe123",
		PhoneNumber: "1234567890",
	}

	mockRepo.On("FindUserByEmail", input.Email).Return(&models.User{}, nil)

	err := userService.SignUp(input)

	assert.EqualError(t, err, models.ErrUserAlreadyExists.Error())
	mockRepo.AssertExpectations(t)

}
