package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		fmt.Println(tokenString)

		if VerifyToken(tokenString) != nil {
			c.JSON(401, gin.H{"error": "invalid or missing token"})
			return
		}

		c.Next()
	}
}
