package main

import (
	"clean-arch/internal/app/config"
	"clean-arch/internal/app/controllers"
	"clean-arch/internal/app/utils"
	"clean-arch/internal/core/database"
	"clean-arch/internal/core/repository"
	"clean-arch/internal/core/services"
	"clean-arch/internal/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()

	log := logger.NewLogrusLogger()

	configEnv := config.ConfigEnv()
	fmt.Println("config env ", configEnv)

	db := database.ConnectDatabase(*configEnv)

	if db == nil {
		log.Error("Failed to connect to database")
		return
	}

	userRepo := repository.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userController := controllers.NewUserController(userService)

	api := gin.Group("/api/v1/users")
	{

		api.POST("/signup", userController.SignUp)
		api.POST("/login", userController.Login)
		api.GET("/profile", utils.AuthMiddleware("user"), userController.GetProfile)
	}

	err := gin.Run(":8080")
	if err != nil {
		log.Error("Failed to start server", err)
	}

}
