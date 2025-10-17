package main

import (
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/evanespen/vanespen.art_2025/internal/albums"
	"github.com/evanespen/vanespen.art_2025/internal/api"
	"github.com/evanespen/vanespen.art_2025/internal/cdn"
	"github.com/evanespen/vanespen.art_2025/internal/pictures"
	"github.com/evanespen/vanespen.art_2025/internal/security"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := api.GetRouter()
	router.Use(cors.Default()) // All origins allowed by default

	router.Use(gin.Recovery())

	adminRouter := api.GetAdminRouter()

	adminRouter.Use(security.AuthMiddleware())

	router.GET("/alive", func(context *gin.Context) {
		context.Status(200)
	})

	pictures.BindRoutes(router, adminRouter)
	albums.BindRoutes(router, adminRouter)
	cdn.BindRoutes(router, adminRouter)
	security.BindRoutes(router, adminRouter)

	err := router.Run(configs.APIHost)
	if err != nil {
		log.Fatal(err)
		return
	}

}
