package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)

var Secret = []byte("your-secret-key")

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) CreateToken(id int, email, role string) (string, error) {
	args := m.Called(id, email, role)
	return args.String(0), args.Error(1)
}

type TokenGenerator interface {
	CreateToken(id int, email, role string) (string, error)
}

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type RealTokenGenerator struct{}

func (r *RealTokenGenerator) CreateToken(id int, email, role string) (string, error) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "The Furnish Store",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret)
}

func AuthMiddleware(requiredRole string, tokenGenerator TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			if requiredRole != "" && claims.Role != requiredRole {
				c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient privileges"})
				c.Abort()
				return
			}
			c.Set("claims", claims)
			c.Set("id", claims.ID)
			c.Set("email", claims.Email)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetClaims(c *gin.Context) (*Claims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, fmt.Errorf("claims not found in context")
	}

	customClaims, ok := claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return customClaims, nil
}
