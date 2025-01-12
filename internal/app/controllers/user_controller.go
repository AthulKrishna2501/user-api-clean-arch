package controllers

import (
	"clean-arch/internal/app/utils"
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/services"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService    services.UserService
	tokenGenerator utils.TokenGenerator
}

func NewUserController(userService services.UserService, tokenGenerator utils.TokenGenerator) *UserController {
	return &UserController{
		userService:    userService,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *UserController) SignUp(ctx *gin.Context) {
	var input models.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := models.ValidateSignup(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.SignUp(&input); err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully!"})
}

func (c *UserController) Login(ctx *gin.Context) {
	var input models.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user, err := c.userService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user.Status == "Blocked" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User is blocked"})
		return
	}

	token, err := c.tokenGenerator.CreateToken(int(user.ID), input.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":           user.ID,
			"user_name":    user.UserName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"status":       user.Status,
			"created_at":   user.CreatedAt.Format(time.RFC3339),
			"updated_at":   user.UpdatedAt.Format(time.RFC3339),
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	customClaims, ok := claims.(*utils.Claims)
	if !ok || customClaims.Email == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetProfile(customClaims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profileResponse := models.UserProfileResponse{
		Name:      user.UserName,
		Email:     user.Email,
		PhnNumber: user.PhoneNumber,
		Status:    user.Status,
	}

	ctx.JSON(http.StatusOK, gin.H{"user": profileResponse})
}
