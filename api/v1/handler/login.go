package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"mes-gin/config"
	"mes-gin/models"
	"mes-gin/utils"
)

var jwtKey = []byte("your_secret_key") // Replace with your actual secret key

// Login handles user login and token generation
func Login(c *gin.Context) {
	var requestBody struct {
		Phone          string `json:"phone" binding:"required"`
		Password       string `json:"password" binding:"required"`
		OrganizationID string `json:"organization_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	baseQuery := config.DB.Preload("Roles").Where("phone = ?", requestBody.Phone)

	if requestBody.OrganizationID != "" {
		baseQuery = baseQuery.Where("organization_id = ?", requestBody.OrganizationID)
	}

	if err := baseQuery.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone"})
		return
	}

	if !utils.CheckPasswordHash(requestBody.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone or password"})
		return
	}

	var userRolesClaims []models.RoleClaims
	for _, role := range user.Roles {
		userRolesClaims = append(userRolesClaims, models.RoleClaims{
			RoleID:    role.ID,
			RoleLabel: role.Label,
			RoleName:  role.Name,
		})
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:         user.ID,
		UserName:       user.Name,
		OrganizationID: user.OrganizationID,
		Roles:          userRolesClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
