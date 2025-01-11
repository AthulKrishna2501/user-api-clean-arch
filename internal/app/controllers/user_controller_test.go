package controllers_test

import (
	"bytes"
	"clean-arch/internal/app/controllers"
	"clean-arch/internal/app/utils"
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/services"
	"clean-arch/internal/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUp(user *models.SignupInput) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Login(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetProfile(userID int) (*models.User, error) {
	args := m.Called(userID)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockTokenGenerator struct{}

func (m *MockTokenGenerator) GenerateToken(userID int, email, role string) (string, error) {
	return "mocked-jwt-token", nil
}

func (m *MockTokenGenerator) CreateToken(id int, email, role string) (string, error) {
	return m.GenerateToken(id, email, role)
}

func TestSignUp(t *testing.T) {
	mockUserService := new(MockUserService)
	mockTokenGenerator := new(MockTokenGenerator)

	userController := controllers.NewUserController(mockUserService, mockTokenGenerator)

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/signup", userController.SignUp)

	input := models.SignupInput{
		UserName:    "JohnDoe",
		Email:       "johndoe@gmail.com",
		Password:    "johndoe123",
		PhoneNumber: "1234567890",
	}

	mockUserService.On("SignUp", &input).Return(nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"message": "User signed up successfully!"}`, rec.Body.String())

	mockUserService.AssertExpectations(t)
}

func TestUserController_SignUp_UserExists(t *testing.T) {
	mockService := new(MockUserService)
	mockTokenGenerator := new(MockTokenGenerator)
	controller := controllers.NewUserController(mockService, mockTokenGenerator)

	router := gin.Default()
	router.POST("/signup", controller.SignUp)

	input := models.SignupInput{
		UserName:    "athul",
		Email:       "athul@gmail.com",
		Password:    "johndoe123",
		PhoneNumber: "1234567890",
	}
	mockService.On("SignUp", &input).Return(models.ErrUserAlreadyExists)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.JSONEq(t, `{"error": "User already exists"}`, rec.Body.String())

	mockService.AssertExpectations(t)
}
func TestLogin_Success(t *testing.T) {

	mockTokenGenerator := new(utils.MockTokenGenerator)

	mockTokenGenerator.On("CreateToken", 1, "johndoe@gmail.com", "user").Return("mocked-jwt-token", nil)

	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService, mockTokenGenerator)

	router := gin.Default()
	router.POST("/login", controller.Login)

	input := models.LoginInput{
		Email:    "johndoe@gmail.com",
		Password: "johndoe123",
	}

	mockUser := &models.User{
		ID:          1,
		UserName:    "JohnDoe",
		Email:       "johndoe@gmail.com",
		PhoneNumber: "1234567890",
		Status:      "Active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("Login", input.Email, input.Password).Return(mockUser, nil)

	mockUser.CreatedAt = mockUser.CreatedAt.Truncate(time.Second)
	mockUser.UpdatedAt = mockUser.UpdatedAt.Truncate(time.Second)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expected := map[string]interface{}{
		"message": "Login successful",
		"token":   "mocked-jwt-token",
		"user": map[string]interface{}{
			"id":           1.0,
			"user_name":    "JohnDoe",
			"email":        "johndoe@gmail.com",
			"phone_number": "1234567890",
			"status":       "Active",
			"created_at":   mockUser.CreatedAt.Format(time.RFC3339),
			"updated_at":   mockUser.UpdatedAt.Format(time.RFC3339),
		},
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error marshaling expected value: %v", err)
	}

	assert.JSONEq(t, string(expectedJSON), rec.Body.String())

	mockService.AssertExpectations(t)
	mockTokenGenerator.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockService := new(MockUserService)
	mockTokenGenerator := new(MockTokenGenerator)
	controller := controllers.NewUserController(mockService, mockTokenGenerator)

	router := gin.Default()
	router.POST("/login", controller.Login)

	input := models.LoginInput{
		Email:    "johndoe@gmail.com",
		Password: "wrongpassword",
	}

	mockService.On("Login", input.Email, input.Password).Return(nil, errors.New("Invalid credentials"))

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	expected := `{"error": "Invalid credentials"}`
	assert.JSONEq(t, expected, rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockService := new(MockUserService)
	mockTokenGenerator := new(MockTokenGenerator)
	controller := controllers.NewUserController(mockService, mockTokenGenerator)

	router := gin.Default()
	router.POST("/login", controller.Login)

	input := models.LoginInput{
		Email:    "notfound@gmail.com",
		Password: "somepassword",
	}

	mockService.On("Login", input.Email, input.Password).Return(nil, errors.New("User not found"))

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	expected := `{"error": "User not found"}`
	assert.JSONEq(t, expected, rec.Body.String())

	mockService.AssertExpectations(t)
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
