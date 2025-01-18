package handlers

import (
	"laundre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			BranchID    uint    `json:"branch_id" binding:"required"`
			Description string  `json:"description" binding:"required"`
			Amount      float64 `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		expense := models.Expense{
			BranchID:    req.BranchID,
			Description: req.Description,
			Amount:      req.Amount,
		}

		if err := db.Create(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "data": expense})
	}
}

func GetAllExpenses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var expenses []models.Expense

		if err := db.Preload("Branch").Find(&expenses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": expenses})
	}
}

func GetExpenseByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Expense ID is required"})
			return
		}

		var expense models.Expense
		if err := db.Preload("Branch").First(&expense, id).Error; err != nil {
			if gorm.ErrRecordNotFound == err {
				c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expense", "details": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": expense})
	}
}

func GetExpensesByBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		branchID := c.Param("branch_id")

		var expenses []models.Expense
		if err := db.Preload("Branch").Where("branch_id = ?", branchID).Find(&expenses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses for the branch", "details": err.Error()})
			return
		}

		if len(expenses) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No expenses found for the specified branch"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": expenses})
	}
}

func UpdateExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var expense models.Expense
		if err := db.First(&expense, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}

		var req struct {
			Description string  `json:"description"`
			Amount      float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		expense.Description = req.Description
		expense.Amount = req.Amount

		if err := db.Save(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully", "data": expense})
	}
}

func DeleteExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result := db.Delete(&models.Expense{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense", "details": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
	}
}
