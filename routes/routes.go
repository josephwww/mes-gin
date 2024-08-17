package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"mes-gin/api/v1/handler"
	"mes-gin/api/v1/middleware"
	_ "mes-gin/docs"
)

// SetupRoutes initializes all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Health check route
	router.GET("/health", handler.HealthCheck)

	//

	// Public routes
	router.POST("/login", handler.Login)
	router.POST("/register", handler.CreateUser)

	apiV1User := router.Group("/api/v1/users")
	apiV1User.Use(middleware.AuthMiddleware())
	{
		// User-related routes
		apiV1User.GET("", handler.GetUsers)
		apiV1User.POST("", handler.CreateUser)
		apiV1User.GET("/:id", handler.GetUser)
		apiV1User.PUT("/:id", handler.UpdateUser)
		apiV1User.DELETE("/:id", handler.DeleteUser)
	}

	apiV1Org := router.Group("/api/v1/orgs")
	apiV1Org.Use(middleware.AuthMiddleware())
	{
		// User-related routes
		apiV1Org.GET("", handler.GetOrganizations)
		apiV1Org.POST("", handler.CreateOrganization)
		apiV1Org.GET("/:id", handler.GetOrganization)
		apiV1Org.PUT("/:id", handler.UpdateOrganization)
		apiV1Org.DELETE("/:id", handler.DeleteOrganization)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Additional groups and routes can be added here
}
