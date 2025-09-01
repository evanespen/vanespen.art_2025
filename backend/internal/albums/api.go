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

	album, err := GetOne(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

type UpdateAlbumPayload struct {
	Pictures []string `json:"pictures"`
}

func UpdateAlbumPictures(c *gin.Context) {
	var payload UpdateAlbumPayload
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	albums, err := Read()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
	}

	fmt.Println(payload)

	var selectedAlbum Album
	var selectedAlbumIndex int

	for index, album := range albums {
		if album.UUID == c.Param("uuid") {
			selectedAlbum = album
			selectedAlbumIndex = index
			break
		}
	}

	if selectedAlbum.Title != "" {
		fmt.Println("album found")
		selectedAlbum.Pictures = payload.Pictures
		albums[selectedAlbumIndex] = selectedAlbum

		Write(albums)
		c.Status(200)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
}

func BindRoutes(engine *gin.Engine, adminGroup *gin.RouterGroup) {
	albumsRouter := engine.Group("/albums")
	albumsRouter.GET("/", GetAllAlbums)
	albumsRouter.GET("/:uuid", GetOneAlbum)
	adminGroup.POST("/albums", PostAlbum)
	adminGroup.PATCH("/albums/:uuid", UpdateAlbumPictures)
}
