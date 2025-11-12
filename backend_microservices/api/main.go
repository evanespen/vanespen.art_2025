// Package main provides a file upload service using MinIO for storage.
//
// This package implements a simple HTTP server that accepts file uploads
// and stores them in a MinIO object storage system. It uses the Gin web
// framework for handling HTTP requests and the MinIO Go SDK for object storage.
package main

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"github.com/vmihailenco/msgpack"

	"vanespen.art-microservices/common/models"
	"vanespen.art-microservices/common/utils"
)

func imageAlreadyExists(mc *minio.Client, objectName string) bool {
	return utils.ObjectExists(mc, "full", objectName) &&
		utils.ObjectExists(mc, "half", objectName) &&
		utils.ObjectExists(mc, "thumb", objectName)
}

func requestGetAll() ([]models.PictureMetadatas, error) {
	nc := utils.NewNatsClient()
	response, err := nc.Request("picture.get_all", nil, 5*time.Second)
	if err != nil {
		log.Println("Got GetAll error:", err)
		return nil, err
	}

	items := []models.PictureMetadatas{}
	err = msgpack.Unmarshal(response.Data, &items)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
		return nil, err
	}

	return items, nil
}

func requestExtract(picture models.Picture) error {
	nc := utils.NewNatsClient()

	// serialize the Picture struct using msgpack
	message, err := msgpack.Marshal(&picture)
	if err != nil {
		log.Println("Failed to marshal picture:", err)
		return err
	}

	msg, err := nc.Request("picture.extract", message, 10*time.Second)
	if err != nil {
		log.Println("Failed to request extract:", err)
		return err
	}

	var response models.ServiceResponse
	err = msgpack.Unmarshal(msg.Data, &response)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
		return err
	}

	log.Println(response)

	return nil
}

// requestResize sends a message to a NATS server with the details of the uploaded picture.
// It uses the msgpack library to serialize the Picture struct and publishes it to the "picture.uploaded" subject.
func requestResize(picture models.Picture) error {

	nc := utils.NewNatsClient()
	defer nc.Close()

	// serialize the Picture struct using msgpack
	message, err := msgpack.Marshal(&picture)
	if err != nil {
		log.Println("Failed to marshal picture:", err)
		return err
	}

	msg, err := nc.Request("picture.resize", message, 10*time.Second)
	if err != nil {
		log.Println("Failed to request picture:", err)
		return err
	}

	var response models.ServiceResponse
	err = msgpack.Unmarshal(msg.Data, &response)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
		return err
	}

	if !response.Success() {
		log.Println("Failed to handle picture:", response.Msg)
		return errors.New(response.Msg)
	}

	return nil
}

// saveFile handles the actual file upload process to MinIO storage
// It takes a multipart.FileHeader as input and uploads the file to the specified bucket
func saveFile(file *multipart.FileHeader) (models.Picture, error) {
	minioClient := utils.NewMinioClient()

	// Open the file for reading
	openedFile, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return models.Picture{}, err
	}
	defer openedFile.Close()

	// Read the file content into a byte slice
	fileBytes := make([]byte, file.Size)
	_, err = openedFile.Read(fileBytes)
	if err != nil {
		fmt.Println("Error reading file content:", err)
		return models.Picture{}, err
	}

	fileExtension := strings.Split(file.Filename, ".")[1]

	bucketName := "full"

	sum := utils.Hash(fileBytes)

	objectName := fmt.Sprintf("%d.%s", sum, fileExtension)

	if imageAlreadyExists(minioClient, objectName) {
		return models.Picture{}, errors.New("file already exists")
	}

	// TODO: Refactor this function to use the utils functions

	if err := utils.PutObject(minioClient, bucketName, objectName, fileBytes); err != nil {
		return models.Picture{}, err
	}

	picture := models.Picture{
		Key:        sum,
		Ext:        fileExtension,
		BytesCount: file.Size,
	}

	log.Printf("File uploaded successfully with name: %s", objectName)

	return picture, nil
}

// uploadHandler handles the HTTP POST request for file uploads
// It verifies the request method, retrieves the file from the form,
// and calls uploadFile to handle the actual upload process
func uploadHandler(c *gin.Context) {
	// Retrieve the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving the file"})
		return
	}

	picture, uploadError := saveFile(file)
	if uploadError != nil {

		if uploadError.Error() == "file already exists" {
			fmt.Println("File already exists:", uploadError)
			c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
			return
		}

		fmt.Println("Error uploading file:", uploadError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file"})
		return
	}
	// c.SSEvent("status", "saving OK: "+picture.Key.String())

	resizeError := requestResize(picture)
	if resizeError != nil {
		fmt.Println("Error resizing:", resizeError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resizing file"})
		return
	}
	// c.SSEvent("status", "resizing OK: "+picture.Key.String())

	extractError := requestExtract(picture)
	if extractError != nil {
		fmt.Println("Error extracting:", extractError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error extracting metadatas"})
		return
	}
	// c.SSEvent("status", "extracting OK: "+picture.Key.String())

	c.JSON(http.StatusOK, gin.H{"message": "File handled", "filename": file.Filename})
}

func getAllHandler(c *gin.Context) {
	items, err := requestGetAll()
	if err != nil {
		fmt.Println("Error getting all items:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting all items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// main initializes the Gin router and starts the HTTP server
// It sets up the POST endpoint for file uploads and starts listening on port 8000
func main() {
	r := gin.Default()

	r.GET("/", getAllHandler)
	r.POST("/", uploadHandler)
	fmt.Println("Starting server at port 8000")
	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
