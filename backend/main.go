package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.GET("/pictures", GetAllPictures)
	router.GET("/pictures/:uuid", GetOnePicture)
	router.POST("/admin/upload", UploadPictures)

	router.GET("/admin/albums", GetAllAlbums)
	router.GET("/admin/albums/:uuid", GetOneAlbum)
	router.POST("/admin/albums", PostAlbum)

	router.Run(":8080")
}
