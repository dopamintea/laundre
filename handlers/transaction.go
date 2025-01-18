package handlers

import (
	"laundre/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionRequest struct {
	CustomerName    string  `json:"customer_name" binding:"required"`
	CustomerPhone   string  `json:"customer_phone" binding:"required"`
	CustomerAddress string  `json:"customer_address" binding:"required"`
	OrderStatus     string  `json:"order_status" binding:"omitempty,oneof=masuk proses urgent done"`
	BranchID        uint    `json:"branch_id" binding:"required"`
	TotalPrice      float64 `json:"total_price" binding:"required"`
	PaymentStatus   string  `json:"payment_status"`
}

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			var customer models.Customer
			if err := tx.Where("name = ? AND phone = ?", req.CustomerName, req.CustomerPhone).
				First(&customer).Error; err != nil {
				if err == gorm.ErrRecordNotFound {

					customer = models.Customer{
						Name:    req.CustomerName,
						Phone:   req.CustomerPhone,
						Address: req.CustomerAddress,
					}
					if err := tx.Create(&customer).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			order := models.Order{
				BranchID:   req.BranchID,
				CustomerID: customer.ID,
				Status:     req.OrderStatus,
			}

			if req.PaymentStatus == "paid" {
				order.Price = req.TotalPrice
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			userID, _ := c.Get("user_id")
			transaction := models.Transaction{
				BranchID:      req.BranchID,
				OrderID:       order.ID,
				UserID:        userID.(uint),
				TotalPrice:    req.TotalPrice,
				PaymentStatus: req.PaymentStatus,
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
	}
}

func GetTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var transactions []models.Transaction
		query := db.Preload("Branch").Preload("Order.Customer").Preload("User")

		if branchID := c.Query("branch_id"); branchID != "" {
			query = query.Where("transactions.branch_id = ?", branchID)
		}
		if status := c.Query("payment_status"); status != "" {
			query = query.Where("transactions.payment_status = ?", status)
		}
		if orderID := c.Query("order_id"); orderID != "" {
			query = query.Where("transactions.order_id = ?", orderID)
		}

		var total int64
		query.Model(&models.Transaction{}).Count(&total)

		if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": transactions,
			"meta": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}

func GetTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var transaction models.Transaction
		if err := db.Preload("Branch").Preload("Order.Customer").Preload("User").
			First(&transaction, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": transaction})
	}
}

func UpdateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var transaction models.Transaction
		if err := db.First(&transaction, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		var req struct {
			PaymentStatus string `json:"payment_status"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.PaymentStatus == "paid" {
			var order models.Order
			if err := db.First(&order, transaction.OrderID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}

			order.Price = transaction.TotalPrice
			if err := db.Save(&order).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order price"})
				return
			}
		}

		transaction.PaymentStatus = req.PaymentStatus

		if err := db.Save(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully", "data": transaction})
	}
}

func DeleteTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result := db.Delete(&models.Transaction{}, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
	}
}
