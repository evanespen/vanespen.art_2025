// Package utils provides utility functions for various operations.
//
// This package contains helper functions that are used across
// different microservices. It includes utilities for MinIO
// object storage and NATS messaging.
package utils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"vanespen.art-microservices/common/models"
)

// NewMinioClient creates a new MinIO client instance.
//
// This function initializes and returns a MinIO client configured
// to connect to the local MinIO server with default credentials.
//
// Returns:
// - *minio.Client: The initialized MinIO client.
func NewMinioClient() *minio.Client {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Cannot create Minio client:", err)
	}

	return minioClient
}

// GetImageBytes retrieves the byte data of an image from MinIO storage.
//
// This function fetches the image data for a given Picture from
// the specified MinIO bucket and returns it as a byte slice.
//
// Parameters:
// - picture: The Picture object for which to retrieve the image data.
//
// Returns:
// - []byte: The byte data of the image.
// - error: An error if the retrieval fails.
func GetImageBytes(picture models.Picture) ([]byte, error) {
	log.Println("Retrieving image:", picture.Key)

	minioClient := NewMinioClient()

	bucketName := "full"
	objectName := fmt.Sprintf("%d.%s", picture.Key, picture.Ext)
	ob, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object:", err)
		return nil, err
	}
	defer ob.Close()

	imageBytes, err := io.ReadAll(ob)
	if err != nil {
		fmt.Println("Error reading object:", err)
		return nil, err
	}

	return imageBytes, nil
}

// GetImage retrieves an image from MinIO storage and decodes it.
//
// This function fetches the image data and decodes it into an
// image.Image object that can be used for further processing.
//
// Parameters:
// - picture: The picture object containing the key and extension.
//
// Returns:
// - image.Image: The decoded image.
// - error: An error if the retrieval or decoding fails.
func GetImage(picture models.Picture) (image.Image, error) {
	imageBytes, err := GetImageBytes(picture)
	if err != nil {
		return nil, err
	}

	// load the bytes to an image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal("Error decoding image:", err)
		return nil, err
	}

	return img, nil
}

// ObjectExists checks if an object exists in a MinIO bucket.
//
// This function verifies the existence of an object in the specified bucket.
// It returns true if the object exists, false otherwise.
//
// Parameters:
// - mc: The MinIO client.
// - bucketName: The name of the bucket.
// - objectName: The name of the object.
//
// Returns:
// - bool: True if the object exists, false otherwise.
func ObjectExists(mc *minio.Client, bucketName string, objectName string) bool {
	_, err := mc.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return true
}

// GetObject retrieves an object from MinIO storage.
//
// This function checks if the specified bucket and object exist,
// then returns the object for reading.
//
// Parameters:
// - mc: The MinIO client.
// - bucketName: The name of the bucket.
// - objectName: The name of the object.
//
// Returns:
// - *minio.Object: The object for reading.
// - error: An error if the object does not exist or if there was an error retrieving it.
func GetObject(mc *minio.Client, bucketName string, objectName string) (*minio.Object, error) {
	// Check if the bucket exists
	exists, err := mc.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucketName)
	}

	if exists := ObjectExists(mc, bucketName, objectName); exists == true {
		return mc.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	} else {
		return nil, fmt.Errorf("object %s does not exist in bucket %s: %v", objectName, bucketName, err)
	}

}

// PutObject uploads an object to MinIO storage.
//
// This function creates a bucket if it doesn't exist and uploads
// the specified data to the specified object in the bucket.
//
// Parameters:
// - mc: The MinIO client.
// - bucketName: The name of the bucket.
// - objectName: The name of the object.
// - data: The data to upload.
//
// Returns:
// - error: An error if the object could not be uploaded.
func PutObject(mc *minio.Client, bucketName string, objectName string, data []byte) error {
	exists, errBucketExists := mc.BucketExists(context.Background(), bucketName)
	if errBucketExists != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", errBucketExists)
	}
	if !exists {
		if err := mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
		}
	}

	info, err := mc.PutObject(
		context.Background(),
		bucketName,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to put object %s in bucket %s: %v", objectName, bucketName, err)
	}

	log.Println("Object put:", info)

	return nil
}
