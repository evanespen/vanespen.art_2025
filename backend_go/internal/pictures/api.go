package pictures

import (
	"fmt"
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func UploadPictures(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		stashImagePath := path.Join(configs.StashDir, file.Filename)
		err := c.SaveUploadedFile(file, stashImagePath)
		if err != nil {
			c.Status(500)
		}

		handleErr := Handle(path.Join(configs.StashDir, file.Filename))

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
	allPictures, err := Read()
	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return
	}
	c.IndentedJSON(http.StatusOK, allPictures)
}

func GetOnePicture(c *gin.Context) {
	uuid := c.Param("uuid")

	allPictures, err := Read()
	if err != nil {
		c.Status(500)
		return
	}

	for _, picture := range allPictures {
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

func BindRoutes(engine *gin.Engine, adminGroup *gin.RouterGroup) {
	picturesRouter := engine.Group("/pictures")
	picturesRouter.GET("/", GetAllPictures)
	picturesRouter.GET("/:uuid", GetOnePicture)
	adminGroup.POST("/pictures", UploadPictures)
}
