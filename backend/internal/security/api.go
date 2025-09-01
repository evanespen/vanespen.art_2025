package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthPayload struct {
	Password string `json:"password"`
}

func HandleAuthenticate(c *gin.Context) {
	var payload AuthPayload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to bind json payload"})
		return
	}

	tokenString, err := Authenticate(payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.Header("Authorization", tokenString)
}

func HandleVerify(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
	}

	err := VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	}

	c.Status(200)
}

func BindRoutes(engine *gin.Engine, adminGroup *gin.RouterGroup) {
	engine.POST("/authenticate", HandleAuthenticate)
	adminGroup.POST("/verify", HandleVerify)
}
