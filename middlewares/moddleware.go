package middlewares

import (
	"github.com/abik1221/Tewanay-Engineering_Intership/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(403, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateAllTokens(clientToken)
		if err != "" {
			c.JSON(403, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)
		c.Set("user_id", claims.User_id)

		c.Next()
	}
}
