package handlers

import (
	"laundre/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRequest struct {
	BranchID   uint   `json:"branch_id" binding:"required"`
	CustomerID uint   `json:"customer_id" binding:"required"`
	Status     string `json:"status" binding:"omitempty,oneof=masuk proses urgent done"`
}

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

		db.Preload("Branch").Preload("Customer").First(&order, order.ID)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Order created successfully",
			"data":    order,
		})
	}
}

func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var orders []models.Order
		query := db.Preload("Branch").Preload("Customer")

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

func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order models.Order
		if err := db.First(&order, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validStatuses := map[string]bool{"masuk": true, "proses": true, "urgent": true, "done": true}
		if !validStatuses[req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		order.Status = req.Status
		order.UpdatedAt = time.Now()

		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Order status updated successfully",
			"data":    order,
		})
	}
}

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
