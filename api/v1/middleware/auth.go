package middleware

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"mes-gin/config"
	"mes-gin/models"
	"net/http"
	"strings"
	"time"
)

// Redis key prefix for user cache
const (
	userCachePrefix         = "current_user_cache:"
	userCacheExpiredSeconds = 1 * time.Hour
	ContextUserKey          = "currentUser"
)

var jwtKey = []byte("your_secret_key") // Replace with your actual secret key

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Try fetching the user from Redis
		var user models.User
		ctx := context.Background()
		cacheKey := userCachePrefix + claims.UserID.String()
		cacheResult, err := config.RDB.Get(ctx, cacheKey).Result()
		if err == redis.Nil {
			// User not found in Redis, fetch from database
			if err := config.DB.First(&user, claims.UserID).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
			// Serialize user object
			userData, _ := json.Marshal(user)

			// Cache user data in Redis with an expiration time
			config.RDB.Set(ctx, cacheKey, userData, userCacheExpiredSeconds) // Cache for 1 hour

			// Store user in context
			c.Set(ContextUserKey, user)
			c.Next()
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		} else {
			// Deserialize user data from Redis
			if err := json.Unmarshal([]byte(cacheResult), &user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user data"})
				c.Abort()
				return
			}
			// Set user in the context for use in handlers
			c.Set("currentUser", user)
			c.Next()
		}
	}
}
