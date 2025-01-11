package handlers

import (
	"laundre/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Request struct for creating/updating customers
type CustomerRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Category string `json:"category" binding:"omitempty,oneof=setia reguler"`
}

// Create a new customer
func CreateCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CustomerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		customer := models.Customer{
			Name:     req.Name,
			Phone:    req.Phone,
			Address:  req.Address,
			Category: req.Category,
		}

		if result := db.Create(&customer); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Customer created successfully",
			"data":    customer,
		})
	}
}

// Get all customers
func GetCustomers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var customers []models.Customer

		query := db.Model(&models.Customer{})
		if err := query.Order("id asc").Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var total int64
		query.Count(&total)

		c.JSON(http.StatusOK, gin.H{
			"data": customers,
			"meta": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}

// Get customer by ID
func GetCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var customer models.Customer
		if err := db.First(&customer, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": customer})
	}
}

// Update customer
func UpdateCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var customer models.Customer
		if err := db.First(&customer, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		var req CustomerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := map[string]interface{}{
			"name":     req.Name,
			"phone":    req.Phone,
			"address":  req.Address,
			"category": req.Category,
		}

		if err := db.Model(&customer).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Customer updated successfully",
			"data":    customer,
		})
	}
}

// Delete customer
func DeleteCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result := db.Delete(&models.Customer{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	}
}
