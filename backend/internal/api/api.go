package main

import (
	"fmt"
	"github.com/evanespen/vanespen.art_2025/internal/albums"
	"github.com/evanespen/vanespen.art_2025/internal/pictures"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func UploadPictures(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		stashImagePath := path.Join(StashDir, file.Filename)
		err := c.SaveUploadedFile(file, stashImagePath)
		if err != nil {
			c.Status(500)
		}

		handleErr := pictures.Handle(path.Join(StashDir, file.Filename))

		if handleErr == nil {
			_ = os.Remove(stashImagePath)
			fmt.Println(fmt.Sprintf("%s process succeed", stashImagePath))
		} else {
			fmt.Println(handleErr)
		}
	}

	c.Status(200)
}

func GetAllPictures(c *gin.Context) {
	pictures, err := pictures.ReadPictures()
	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return
	}
	c.IndentedJSON(http.StatusOK, pictures)
}

func GetOnePicture(c *gin.Context) {
	uuid := c.Param("uuid")

	pictures, err := pictures.ReadPictures()
	if err != nil {
		c.Status(500)
		return
	}

	for _, picture := range pictures {
		if picture.UUID == uuid {
			c.IndentedJSON(http.StatusOK, picture)
			return
		}
	}

	c.Status(404)
	return
}

// TODO: DELETE
// TODO: PUT (favourite)

func PostAlbum(c *gin.Context) {
	newAlbum := NewAlbum()
	if err := c.BindJSON(newAlbum); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	albums.AppendAlbum(*newAlbum)

	c.Status(201)
}

func GetAllAlbums(c *gin.Context) {
	albums, err := albums.ReadAlbums()
	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func GetOneAlbum(c *gin.Context) {
	uuid := c.Param("uuid")

	albums, err := albums.ReadAlbums()
	if err != nil {
		c.Status(500)
		return
	}

	for _, album := range albums {
		if album.UUID == uuid {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.Status(404)
	return
}
