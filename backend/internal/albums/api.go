package albums

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostAlbum(c *gin.Context) {
	newAlbum := NewAlbum()
	if err := c.BindJSON(newAlbum); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	Append(*newAlbum)

	c.Status(201)
}

func GetAllAlbums(c *gin.Context) {
	allAlbums, err := Read()
	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return
	}
	c.IndentedJSON(http.StatusOK, allAlbums)
}

func GetOneAlbum(c *gin.Context) {
	uuid := c.Param("uuid")

	allAlbums, err := Read()
	if err != nil {
		c.Status(500)
		return
	}

	for _, album := range allAlbums {
		if album.UUID == uuid {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	c.Status(404)
	return
}

func BindRoutes(engine *gin.Engine, adminGroup *gin.RouterGroup) {
	albumsRouter := engine.Group("/albums")
	albumsRouter.GET("/", GetAllAlbums)
	albumsRouter.GET("/:uuid", GetOneAlbum)
	adminGroup.POST("/albums", PostAlbum)
}
