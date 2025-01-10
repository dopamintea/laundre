package handlers

import (
	"laundre/models"
	"laundre/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Cari user
		var user models.User
		if err := db.Where("username = ? AND status = ?", req.Username, "active").First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Verif password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Gen token
		token, err := utils.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Update last login
		now := time.Now()
		user.LastLogin = &now
		db.Save(&user)

		// Bikin login log
		log := models.Log{
			UserID:    user.ID,
			IPAddress: c.ClientIP(),
		}
		db.Create(&log)

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
				"branchId": user.BranchID,
			},
		})
	}
}

// handlers/auth.go
func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		token, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found"})
			return
		}

		// Get token claims to get expiration time
		claims, err := utils.ValidateToken(token.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		// Add token to blacklist
		blacklistedToken := models.TokenBlacklist{
			Token:     token.(string),
			UserID:    userID.(uint),
			ExpiresAt: claims.ExpiresAt.Time,
		}

		if err := db.Create(&blacklistedToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}

		// Create logout log
		logEntry := models.Log{
			UserID:    userID.(uint),
			IPAddress: c.ClientIP(),
		}
		db.Create(&logEntry)

		// Cleanup expired tokens
		go models.CleanupBlacklist(db)

		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}
