package handlers

import (
	"laundre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin staf"`
	BranchID *uint  `json:"branch_id"`
	Status   string `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role" binding:"omitempty,oneof=admin staf"`
	BranchID *uint  `json:"branch_id"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Role == "staf" && req.BranchID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch ID is required for staff users"})
			return
		}

		if req.BranchID != nil {
			var branch models.Branch
			if err := db.First(&branch, req.BranchID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
				return
			}
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Username: req.Username,
			Password: string(hashedPassword),
			Role:     req.Role,
			BranchID: req.BranchID,
			Status:   req.Status,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		user.Password = ""
		c.JSON(http.StatusCreated, user)
	}
}

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Select("id, username, role, branch_id, status, last_login").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.Select("id, username, role, branch_id, status, last_login").First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Role == "staf" && req.BranchID == nil && user.BranchID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch ID is required for staff users"})
			return
		}

		if req.BranchID != nil {
			var branch models.Branch
			if err := db.First(&branch, req.BranchID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
				return
			}
		}

		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			user.Password = string(hashedPassword)
		}
		if req.Role != "" {
			user.Role = req.Role
		}
		if req.Status != "" {
			user.Status = req.Status
		}
		if req.BranchID != nil {
			user.BranchID = req.BranchID
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		user.Password = ""
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
