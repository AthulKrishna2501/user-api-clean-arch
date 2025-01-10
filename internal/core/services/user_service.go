package services

import (
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(user *models.SignupInput) error
	Login(email, password string) (*models.User, error)
	GetProfile(userID int) (*models.User, error)
}

type UserServiceImpl struct {
	userRepo repository.UserRespository
}

func NewUserService(userRepo repository.UserRespository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) SignUp(user *models.SignupInput) error {
	exists, _ := s.userRepo.FindUserByEmail(user.Email)
	if exists != nil {
		return errors.New(models.ErrUserAlreadyExists)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	newUser := &models.User{
		UserName:    user.UserName,
		Email:       user.Email,
		Password:    string(hashedPassword),
		PhoneNumber: user.PhoneNumber,
		Status:      "Active",
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}

	return nil

}

func (s *UserServiceImpl) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New(models.ErrUserDoesNotExist)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil

}

func (s *UserServiceImpl) GetProfile(userID int) (*models.User, error) {
	return s.userRepo.FindUserByID(userID)
}
