package repository_test

import (
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := repository.NewUserRepository(gormDB)

	user := &models.User{
		UserName:    "JohnDoe",
		Email:       "johndoe@gmail.com",
		Password:    "johndoe123",
		PhoneNumber: "1234567890",
		Status:      "Active",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			user.UserName,
			user.Email,
			user.Password,
			user.PhoneNumber,
			user.Status,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestFindUserByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := repository.NewUserRepository(gormDB)

	email := "johndoe@gmail.com"
	mockUser := models.User{
		UserName: "JohnDoe",
		Email:    email,
		Status:   "Active",
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"."deleted_at" IS NULL`).WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_name", "email", "status"}).
			AddRow(mockUser.UserName, mockUser.Email, mockUser.Status))

	user, err := repo.FindUserByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, mockUser.UserName, user.UserName)
	assert.Equal(t, mockUser.Email, user.Email)
	assert.Equal(t, mockUser.Status, user.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}
