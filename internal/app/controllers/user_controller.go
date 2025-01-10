package controllers

import (
	"clean-arch/internal/app/utils"
	"clean-arch/internal/core/models"
	"clean-arch/internal/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) SignUp(ctx *gin.Context) {
	var input models.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})

		return
	}

	if err := models.ValidateSignup(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.SignUp(&input); err != nil {
		if err.Error() == models.ErrUserAlreadyExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": models.MsgSignupSuccessful})

}

func (c *UserController) Login(ctx *gin.Context) {
	var input models.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidInput})
		return
	}

	if input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidInput})
		return
	}

	user, err := c.userService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user.Status == "Blocked" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrUserBlocked})
		return
	}

	token, err := utils.CreateToken(1, input.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": models.MsgLoginSuccessful,
		"user":    user,
		"token":   token,
	})

}

func (c *UserController) GetProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidID})
		return
	}

	customClaims, ok := claims.(*utils.Claims)
	if !ok || customClaims.Email == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidID})
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
