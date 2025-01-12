package handlers

import (
	"laundre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateBranchRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
}

type UpdateBranchRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func CreateBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateBranchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		branch := models.Branch{
			Name:    req.Name,
			Address: req.Address,
			Phone:   req.Phone,
		}

		if err := db.Create(&branch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
			return
		}

		c.JSON(http.StatusCreated, branch)
	}
}

func GetBranches(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var branches []models.Branch
		if err := db.Find(&branches).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branches"})
			return
		}

		c.JSON(http.StatusOK, branches)
	}
}

func GetBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var branch models.Branch
		if err := db.First(&branch, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}

		c.JSON(http.StatusOK, branch)
	}
}

func UpdateBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var branch models.Branch
		if err := db.First(&branch, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}

		var req UpdateBranchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Name != "" {
			branch.Name = req.Name
		}
		if req.Address != "" {
			branch.Address = req.Address
		}
		if req.Phone != "" {
			branch.Phone = req.Phone
		}

		if err := db.Save(&branch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch"})
			return
		}

		c.JSON(http.StatusOK, branch)
	}
}

func DeleteBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var branch models.Branch
		if err := db.First(&branch, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}

		if err := db.Delete(&branch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
	}
}
