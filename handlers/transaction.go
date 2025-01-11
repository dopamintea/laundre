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
	OrderStatus     string  `json:"order_status" binding:"omitempty,oneof=masuk proses urgent"`
	BranchID        uint    `json:"branch_id" binding:"required"`
	TotalPrice      float64 `json:"total_price" binding:"required"`
	PaymentStatus   string  `json:"payment_status"`
}

// Create transaction
func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Begin database transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			// Step 1: Find or create customer
			var customer models.Customer
			if err := tx.Where("name = ? AND phone = ?", req.CustomerName, req.CustomerPhone).
				First(&customer).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create new customer
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

			// Step 2: Create order
			order := models.Order{
				BranchID:   req.BranchID,
				CustomerID: customer.ID,
				Status:     req.OrderStatus,
			}
			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			// Step 3: Create transaction
			userID, _ := c.Get("user_id") // Assuming user_id is set in context by authentication middleware
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

			// Return nil to commit the transaction
			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
	}
}

// Get all transactions
func GetTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var transactions []models.Transaction
		query := db.Preload("Branch").Preload("Order.Customer").Preload("User")

		// Apply filters
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

// Get a transaction by ID
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

// Update a transaction
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

		transaction.PaymentStatus = req.PaymentStatus

		if err := db.Save(&transaction).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully", "data": transaction})
	}
}

// Delete a transaction
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
