package main

import (
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/evanespen/vanespen.art_2025/internal/albums"
	"github.com/evanespen/vanespen.art_2025/internal/api"
	"github.com/evanespen/vanespen.art_2025/internal/pictures"
	"github.com/gin-gonic/gin"
)

func main() {
	router := api.GetRouter()
	adminRouter := api.GetAdminRouter()

	router.GET("/alive", func(context *gin.Context) {
		context.Status(200)
	})

	pictures.BindRoutes(router, adminRouter)
	albums.BindRoutes(router, adminRouter)

	router.Run(configs.APIHost)
}
