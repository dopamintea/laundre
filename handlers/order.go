package handlers

import (
	"laundre/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Request struct for creating/updating orders
type OrderRequest struct {
	BranchID   uint   `json:"branch_id" binding:"required"`
	CustomerID uint   `json:"customer_id" binding:"required"`
	Status     string `json:"status" binding:"omitempty,oneof=masuk proses urgent done"`
}

// Create order
func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order := models.Order{
			BranchID:   req.BranchID,
			CustomerID: req.CustomerID,
			Status:     req.Status,
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Fetch complete order with relationships
		db.Preload("Branch").Preload("Customer").First(&order, order.ID)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Order created successfully",
			"data":    order,
		})
	}
}

// Get all orders
func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var orders []models.Order
		query := db.Preload("Branch").Preload("Customer")

		// Apply filters
		if branchID := c.Query("branch_id"); branchID != "" {
			query = query.Where("orders.branch_id = ?", branchID)
		}
		if customerID := c.Query("customer_id"); customerID != "" {
			query = query.Where("orders.customer_id = ?", customerID)
		}
		if status := c.Query("status"); status != "" {
			query = query.Where("orders.status = ?", status)
		}

		var total int64
		query.Model(&models.Order{}).Count(&total)

		if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": orders,
			"meta": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}

// Get order by ID
func GetOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order models.Order
		if err := db.Preload("Branch").Preload("Customer").First(&order, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": order})
	}
}

// Update order
func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order models.Order
		if err := db.First(&order, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		var req OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Only update allowed fields
		updates := map[string]interface{}{
			"status": req.Status,
		}

		if err := db.Model(&order).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Fetch updated order with relationships
		db.Preload("Branch").Preload("Customer").First(&order, id)

		c.JSON(http.StatusOK, gin.H{
			"message": "Order updated successfully",
			"data":    order,
		})
	}
}

// Delete order
func DeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result := db.Delete(&models.Order{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
}
