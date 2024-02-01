package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test/pkg/utils"
)

// BearerAuthMiddleware is a middleware function to check for Bearer token in the request headers
func BearerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Bearer token from the Authorization header
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		Claims, boolean := isValidToken(token)
		if boolean {
			fmt.Println("Claims:")
			for key, value := range Claims {
				fmt.Printf("%s: %v\n", key, value)
			}
			// Store the claims in the Gin context
			c.Set("claims", Claims)
			c.Next() // Token is valid, continue to the next handler
		} else {
			// Token is invalid, return an unauthorized response
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}

func isValidToken(token string) (map[string]interface{}, bool) {
	Claims, IsValid, _ := utils.JWTDecodeToken(token)
	return Claims, IsValid
}
