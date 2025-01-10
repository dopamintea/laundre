// middleware/auth.go
package middleware

import (
	"laundre/models"
	"laundre/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		// Check if token is blacklisted
		var blacklistedToken models.TokenBlacklist
		if err := db.Where("token = ?", tokenString).First(&blacklistedToken).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store token string and claims in context
		c.Set("token", tokenString) // Store the raw token string
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("branch_id", claims.BranchID)

		c.Next()
	}
}
