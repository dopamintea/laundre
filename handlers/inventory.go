package handlers

import (
	"laundre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateInventory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			BranchID uint   `json:"branch_id" binding:"required"`
			Name     string `json:"name" binding:"required"`
			Stock    int    `json:"stock" binding:"required,min=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var branch models.Branch
		if err := db.First(&branch, req.BranchID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}

		inventory := models.Inventory{
			BranchID: req.BranchID,
			Name:     req.Name,
			Stock:    req.Stock,
		}

		if err := db.Create(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Inventory created successfully", "data": inventory})
	}
}

func GetAllInventories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inventories []models.Inventory

		if err := db.Preload("Branch").Find(&inventories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventories", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": inventories})
	}
}

func GetInventoryByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var inventory models.Inventory
		if err := db.Preload("Branch").First(&inventory, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": inventory})
	}
}

func UpdateInventory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var inventory models.Inventory
		if err := db.First(&inventory, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
			return
		}

		var req struct {
			Name  string `json:"name" binding:"omitempty,max=100"`
			Stock *int   `json:"stock" binding:"omitempty,min=0"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		if req.Name != "" {
			inventory.Name = req.Name
		}
		if req.Stock != nil {
			inventory.Stock = *req.Stock
		}

		if err := db.Save(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully", "data": inventory})
	}
}

func DeleteInventory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result := db.Delete(&models.Inventory{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory", "details": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
	}
}

func GetInventoryByBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			BranchID uint `json:"branch_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		var inventories []models.Inventory

		if err := db.Preload("Branch").Where("branch_id = ?", req.BranchID).Find(&inventories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventories for the branch", "details": err.Error()})
			return
		}

		if len(inventories) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No inventory found for the specified branch"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": inventories})
	}
}
