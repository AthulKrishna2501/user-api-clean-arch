package services_test

import (
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/services"
	"clean-arch/internal/mocks"
	"testing"
	"time"

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
func TestGetProfile_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	service := services.NewUserService(mockRepo)

	mockUser := &models.User{
		ID:          1,
		UserName:    "JohnDoe",
		Email:       "johndoe@gmail.com",
		PhoneNumber: "1234567890",
		Status:      "Active",
		CreatedAt:   time.Now().Truncate(time.Second),
		UpdatedAt:   time.Now().Truncate(time.Second),
	}

	mockRepo.On("FindUserByID", 1).Return(mockUser, nil)

	result, err := service.GetProfile(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, mockUser.ID, result.ID)
	assert.Equal(t, mockUser.UserName, result.UserName)
	assert.Equal(t, mockUser.Email, result.Email)
	assert.Equal(t, mockUser.PhoneNumber, result.PhoneNumber)
	assert.Equal(t, mockUser.Status, result.Status)
	assert.Equal(t, mockUser.CreatedAt, result.CreatedAt)
	assert.Equal(t, mockUser.UpdatedAt, result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}
