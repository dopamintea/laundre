package handlers

import (
	"database/sql"
	"laundre/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetGrossProfit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalGrossProfit sql.NullFloat64

		if err := db.Model(&models.Order{}).Where("status = ?", "done").Select("sum(price)").Scan(&totalGrossProfit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate gross profit from orders", "details": err.Error()})
			return
		}

		if !totalGrossProfit.Valid {
			totalGrossProfit.Float64 = 0.0
		}

		c.JSON(http.StatusOK, gin.H{"gross_profit": totalGrossProfit.Float64})
	}
}

func GetProfit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var grossProfit sql.NullFloat64
		var totalExpenses sql.NullFloat64

		if err := db.Model(&models.Order{}).Where("status = ?", "done").
			Select("sum(price)").Scan(&grossProfit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate gross profit from orders", "details": err.Error()})
			return
		}

		if !grossProfit.Valid {
			grossProfit.Float64 = 0.0
		}

		if err := db.Model(&models.Expense{}).Select("sum(amount)").Scan(&totalExpenses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total expenses", "details": err.Error()})
			return
		}

		if !totalExpenses.Valid {
			totalExpenses.Float64 = 0.0
		}

		netProfit := grossProfit.Float64 - totalExpenses.Float64

		c.JSON(http.StatusOK, gin.H{
			"gross_profit":   grossProfit.Float64,
			"total_expenses": totalExpenses.Float64,
			"net_profit":     netProfit,
		})
	}
}

func GetGrossProfitByBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		branchID := c.Param("branch_id")

		if branchID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch ID is required"})
			return
		}

		var grossProfit sql.NullFloat64

		if err := db.Model(&models.Order{}).Where("status = ? AND branch_id = ?", "done", branchID).
			Select("sum(price)").Scan(&grossProfit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate gross profit from orders", "details": err.Error()})
			return
		}

		if !grossProfit.Valid {
			grossProfit.Float64 = 0.0
		}

		c.JSON(http.StatusOK, gin.H{
			"branch_id":    branchID,
			"gross_profit": grossProfit.Float64,
		})
	}
}

func GetProfitByBranch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		branchID := c.Param("branch_id")

		if branchID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch ID is required"})
			return
		}

		var grossProfit sql.NullFloat64
		var totalExpenses sql.NullFloat64

		if err := db.Model(&models.Order{}).Where("status = ? AND branch_id = ?", "done", branchID).
			Select("sum(price)").Scan(&grossProfit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate gross profit from orders", "details": err.Error()})
			return
		}

		if !grossProfit.Valid {
			grossProfit.Float64 = 0.0
		}

		if err := db.Model(&models.Expense{}).Where("branch_id = ?", branchID).
			Select("sum(amount)").Scan(&totalExpenses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total expenses", "details": err.Error()})
			return
		}

		if !totalExpenses.Valid {
			totalExpenses.Float64 = 0.0
		}

		netProfit := grossProfit.Float64 - totalExpenses.Float64

		c.JSON(http.StatusOK, gin.H{
			"branch_id":      branchID,
			"gross_profit":   grossProfit.Float64,
			"total_expenses": totalExpenses.Float64,
			"net_profit":     netProfit,
		})
	}
}
