package middlewares

import (
	"net/http"
	"strings"
	"abc.com/calc/models"
	"abc.com/calc/utils"
	"github.com/gin-gonic/gin"
)


func Authenticate(context *gin.Context) {
	// Get the Authorization header
	authHeader := context.Request.Header.Get("Authorization")

	// Check if the header is empty
	if authHeader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not logged in "})
		return
	}

	// Split the header to check for "Bearer" prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	// Extract the token
	token := parts[1]

	// Verify the token
	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	// Store the user ID in the context
	context.Set("userId", userId)
	context.Next()
}


func AuthorizeAdmin(context *gin.Context) {
	// Retrieve user ID from context (set by Authenticate middleware)
	userID, exists := context.Get("userId")
	if !exists {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please log in."})
		return
	}

	// Convert userID to int64
	userIDInt, ok := userID.(int64)
	if !ok {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format."})
		return
	}

	// Get the role of the user
	role, err := models.GetRoleById(userIDInt)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user role."})
		return
	}




	// Check if the role is "administrator"
	if role != "administrator" {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden. Administrator access required."})
		return
	}

	// Proceed to the next middleware or handler
	context.Next()
}
