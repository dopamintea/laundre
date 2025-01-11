package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOrStaffBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Admin has full access
		if role == "admin" {
			c.Next()
			return
		}

		// For staff, check branch access
		if role == "staf" {
			userBranchID, exists := c.Get("branch_id")
			if !exists || userBranchID == nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "No branch assigned"})
				c.Abort()
				return
			}

			// Get branch_id from request (same logic as BranchAccess middleware)
			requestedBranchID := c.Param("branch_id")
			if requestedBranchID == "" {
				requestedBranchID = c.Query("branch_id")
				if requestedBranchID == "" {
					if c.Request.Method != "GET" {
						var bodyMap map[string]interface{}
						if err := c.ShouldBindJSON(&bodyMap); err == nil {
							if branchID, exists := bodyMap["branch_id"]; exists {
								// Mengonversi float64 ke string
								requestedBranchID = fmt.Sprintf("%.0f", branchID.(float64))
							}
						}
						// Mendapatkan kembali body request
						body, err := c.Request.GetBody()
						if err == nil {
							c.Request.Body = body
						} else {
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request body"})
							return
						}
					}
				}
			}

			if requestedBranchID != "" && requestedBranchID != string(userBranchID.(uint)) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for this branch"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
