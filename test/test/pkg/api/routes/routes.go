package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"test/pkg/api/handlers"
	"test/pkg/api/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Use the default configuration for CORS middleware
	router.Use(cors.Default())

	authorized := router.Group("/")
	authorized.Use(middleware.BearerAuthMiddleware())

	// Users Route
	router.GET("/users", handlers.GetUsersHandler)
	router.GET("/user/:id", handlers.GetUserHandler)
	router.DELETE("/user/:id", handlers.DeleteUserHandler)
	router.POST("/user/signup", handlers.AddUserHandler)
	router.POST("/user/signin", handlers.LoginUserHandler)
	router.POST("/refresh-token", handlers.RefreshTokenHandler)
	router.PUT("/user/:id", handlers.UpdateUserHandler)

	// Organizations Route
	authorized.GET("/organizations", handlers.GetOrganizationsHandler)
	authorized.GET("/organization/:id", handlers.GetOrganizationHandler)
	authorized.POST("/organization", handlers.AddOrganizationHandler)
	authorized.POST("/organization/:organization_id/invite", handlers.InviteOrganizationHandler)
	authorized.PUT("/organization/:organization_id", handlers.UpdateOrganizationHandler)
	authorized.DELETE("/organization/:organization_id", handlers.DeleteOrganizationHandler)
	//Test
	router.GET("/GoApi", handlers.TestHandler)
}
