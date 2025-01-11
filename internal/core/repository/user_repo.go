package repository

import (
	"clean-arch/internal/core/models"
	"errors"

	"gorm.io/gorm"
)

type UserStorage struct {
	DB *gorm.DB
}

type UserRespository interface {
	FindUserByEmail(string) (*models.User, error)
	FindUserByID(int) (*models.User, error)
	CreateUser(*models.User) error
	// UpdateUser(*models.User) error
}

func NewUserRepository(db *gorm.DB) *UserStorage {
	return &UserStorage{
		DB: db,
	}
}

func (repo *UserStorage) CreateUser(user *models.User) error {
	if err := repo.DB.Create(user).Error; err != nil {
		return errors.New("failed to create user: " + err.Error())
	}

	return nil
}

func (repo *UserStorage) FindUser(field string, value interface{}) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where(field+" = ?", value).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user: " + err.Error())
	}
	return &user, nil
}

func (repo *UserStorage) FindUserByEmail(email string) (*models.User, error) {
	return repo.FindUser("email", email)
}

func (repo *UserStorage) FindUserByID(userID int) (*models.User, error) {
	return repo.FindUser("id", userID)
}

// func (repo *UserStorage) UpdateUser(user *models.User) error {
// 	if err := repo.DB.Save(user).Error; err != nil {
// 		return errors.New("failed to update user:" + err.Error())
// 	}

// 	return nil
// }
