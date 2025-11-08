// Package main provides a file upload service using MinIO for storage.
//
// This package implements a simple HTTP server that accepts file uploads
// and stores them in a MinIO object storage system. It uses the Gin web
// framework for handling HTTP requests and the MinIO Go SDK for object storage.
package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/nats-io/nats.go"

	"github.com/vmihailenco/msgpack"
)

// Picture represents a picture stored in MinIO with its metadata
type Picture struct {
	Key        uuid.UUID
	Ext        string
	BytesCount int64
}

// notify sends a message to a NATS server with the details of the uploaded picture.
// It uses the msgpack library to serialize the Picture struct and publishes it to the "picture.uploaded" subject.
func notify(picture Picture) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// serialize the Picture struct using msgpack
	message, err := msgpack.Marshal(&picture)
	if err != nil {
		log.Fatal("Failed to marshal picture:", err)
	}

	// notify the NATS server about the uploaded picture
	err = nc.Publish("picture.uploaded", message)
	if err != nil {
		log.Fatal("Failed to publish picture:", err)
	}
}

// uploadFile handles the actual file upload process to MinIO storage
// It takes a multipart.FileHeader as input and uploads the file to the specified bucket
func uploadFile(file *multipart.FileHeader) {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println("Error creating Minio client:", err)
		return
	}

	// Extract file extension from the filename
	fileExtension := strings.Split(file.Filename, ".")[1]

	// Define the bucket name where files will be stored
	bucketName := "full"
	// Generate a new UUID for the object name
	objectUuid := uuid.New()
	objectName := fmt.Sprintf("%s.%s", objectUuid.String(), fileExtension)

	// Open the file for reading
	openedFile, err := file.Open()
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer openedFile.Close()

	// Check if the bucket exists, if not, create it
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println("Error checking bucket existence:", err)
		return
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println("Error creating bucket:", err)
			return
		}
	}

	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		openedFile,
		file.Size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	picture := Picture{
		Key:        objectUuid,
		Ext:        fileExtension,
		BytesCount: file.Size,
	}

	notify(picture)

	log.Printf("File uploaded successfully with name: %s", objectName)
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

	uploadFile(file)

	c.JSON(http.StatusOK, gin.H{"message": "File handled", "filename": file.Filename})
}

// main initializes the Gin router and starts the HTTP server
// It sets up the POST endpoint for file uploads and starts listening on port 8000
func main() {
	r := gin.Default()

	r.POST("/", uploadHandler)
	fmt.Println("Starting server at port 8000")
	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
