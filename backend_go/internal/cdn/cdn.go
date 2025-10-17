package cdn

import (
	"errors"
	"fmt"
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

type ImageSize string

const (
	Full  ImageSize = "full"
	Half  ImageSize = "half"
	Thumb ImageSize = "thumb"
	Tiny  ImageSize = "tiny"
)

func GetImagePath(uuid string, ext string, imageSize ImageSize) (string, error) {
	var imagePath string
	switch imageSize {
	case Full:
		imagePath = path.Join(configs.FullResDir, uuid+ext)
	case Half:

		imagePath = path.Join(configs.HalfResDir, uuid+ext)
	case Thumb:
		imagePath = path.Join(configs.ThumbResDir, uuid+ext)
	case Tiny:
		imagePath = path.Join(configs.TinyResDir, uuid+ext)
	}

	if _, err := os.Stat(imagePath); err != nil {
		return "", errors.New("file not found")
	}

	return imagePath, nil
}

func HandleGetImage(c *gin.Context) {
	uuid := c.Param("uuid")
	ext := c.DefaultQuery("ext", ".jpg")
	sizeRaw := c.DefaultQuery("size", "thumb")

	var size ImageSize
	switch sizeRaw {
	case "thumb":
		size = Thumb
	case "half":
		size = Half
	case "full":
		size = Full
	case "tiny":
		size = Tiny
	default:
		size = Thumb
	}

	fmt.Println(uuid, ext, size)
	imagePath, err := GetImagePath(uuid, ext, size)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.File(imagePath)
}

func BindRoutes(engine *gin.Engine, adminGroup *gin.RouterGroup) {
	cdnRouter := engine.Group("/cdn")
	cdnRouter.GET("/:uuid", HandleGetImage)
}
