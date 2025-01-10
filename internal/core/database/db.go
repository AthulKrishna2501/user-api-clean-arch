package database

import (
	"clean-arch/internal/app/config"
	"clean-arch/internal/core/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(env config.Env) *gorm.DB {
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		env.DBHOST,
		env.DBUSER,
		env.DBPASSWORD,
		env.DBNAME,
		env.DBPORT,
		env.SSLMODE,
	)

	fmt.Println("DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed", err)
		return nil
	}

	err = AutoMigrate(db)
	if err != nil {
		log.Fatal("error in automigration", err)
		return nil
	}
	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.TempUser{},
	)
}
