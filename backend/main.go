package main

import (
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/evanespen/vanespen.art_2025/internal/albums"
	"github.com/evanespen/vanespen.art_2025/internal/api"
	"github.com/evanespen/vanespen.art_2025/internal/pictures"
	"github.com/evanespen/vanespen.art_2025/internal/security"
	"github.com/gin-gonic/gin"
)

func main() {
	router := api.GetRouter()

	router.Use(gin.Recovery())

	adminRouter := api.GetAdminRouter()

	adminRouter.Use(security.AuthMiddleware())

	router.GET("/alive", func(context *gin.Context) {
		context.Status(200)
	})

	pictures.BindRoutes(router, adminRouter)
	albums.BindRoutes(router, adminRouter)
	security.BindRoutes(router, adminRouter)

	router.Run(configs.APIHost)

}
