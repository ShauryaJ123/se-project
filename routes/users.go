package routes

import (
	"log"
	"net/http"

	"abc.com/calc/models"
	"abc.com/calc/utils"
	"github.com/gin-gonic/gin"
)


func signup(context *gin.Context) {
	var user models.User

	// Attempt to bind the incoming JSON to the User struct
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input. Please provide valid user details."})
		return
	}
	// Additional validation check for mandatory fields (example: email, password)
	if user.Email == "" || user.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required fields."})
		return
	}
	// Attempt to save the user
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user. Please try again later."})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}




func login(context *gin.Context) {
	var user models.User2
	// Attempt to bind the incoming JSON to the User struct
	err := context.ShouldBindJSON(&user)
	log.Println(user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input. Please provide valid login details."})
		log.Println(err)
		return
	}
	// Ensure both email and password are provided
	if user.Email == "" || user.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required fields."})
		return
	}
	// Check the user's credentials
	err = user.CheckCredentials2()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password. Please try again."})
		return
	}
	// Generate a JWT token upon successful authentication
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while generating the authentication token. Please try again later."})
		return
	}
	// Respond with the generated token
	context.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully.",
		"token":   token,
	})
}

