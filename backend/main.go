package main

import (
	"github.com/evanespen/vanespen.art_2025/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/pictures", api.GetAllPictures)
	router.GET("/pictures/:uuid", api.GetOnePicture)
	router.POST("/admin/upload", api.UploadPictures)

	router.GET("/albums", api.GetAllAlbums)
	router.GET("/albums/:uuid", api.GetOneAlbum)
	router.POST("/admin/albums", api.PostAlbum)

	router.Run(":8080")
}
