package api

import "github.com/gin-gonic/gin"

var router = gin.Default()

func GetRouter() *gin.Engine {
	return router
}

func GetAdminRouter() *gin.RouterGroup {
	return GetRouter().Group("/admin")
}
