package handler

import (
	"mes-gin/api/v1/request"
	"net/http"

	"github.com/gin-gonic/gin"

	"mes-gin/config"
	"mes-gin/models"
	"mes-gin/utils"
)

// GetUsers returns a list of users
func GetUsers(c *gin.Context) {
	var listParams request.Pagination
	if err := c.BindQuery(&listParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	users, users_count, err := models.GetAllUsers(config.DB, listParams)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": users_count, "data": users})
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user models.User

	// Parse the request body into the User model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := models.CreateUser(config.DB, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Respond with the created user details
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUser returns a specific user by ID
func GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := models.GetUserByID(config.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a specific user by ID
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	// Placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"id":      userID,
	})
}

// DeleteUser deletes a specific user by ID
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	// Placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
		"id":      userID,
	})
}

// CurrentUser returns the currently authenticated user
func CurrentUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
