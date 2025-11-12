// Package main provides database operations for the application.
//
// This package handles all database-related functionality, including
// connecting to the database, performing CRUD operations, and managing
// data storage. It uses MinIO for object storage and parquet format
// for efficient data serialization.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/parquet-go/parquet-go"

	// minioPq "github.com/xitongsys/parquet-go-source/minio"
	// "github.com/xitongsys/parquet-go/parquet"
	// "github.com/xitongsys/parquet-go/reader"
	// "github.com/xitongsys/parquet-go/writer"
	"vanespen.art-microservices/common/models"
	"vanespen.art-microservices/common/utils"
)

// Insert adds a new picture metadata record to the database.
//
// This function appends the new item to the existing records (if any)
// and writes the updated list to the database.
//
// Parameters:
//   - mc: The MinIO client instance
//   - item: The PictureMetadatas struct containing the metadata to insert
//
// Returns:
//   - error: An error if the insertion or writing process fails
func Insert(mc *minio.Client, item models.PictureMetadatas) error {
	pictures := []models.PictureMetadatas{}

	if utils.ObjectExists(mc, "database", "pictures.parquet") == true {
		_pictures, err := GetAll(mc)
		if err != nil {
			log.Println("Failed to get pictures:", err)
			return err
		}
		pictures = append(_pictures, item)
	} else {
		pictures = append(pictures, item)
	}

	return Write(mc, pictures)
}

// Write saves picture metadata to the database.
//
// This function writes a slice of PictureMetadatas to a parquet file
// in the MinIO database bucket.
//
// Parameters:
//   - mc: The MinIO client instance
//   - items: A slice of PictureMetadatas containing the metadata to write
//
// Returns:
//   - error: An error if the writing process fails
func Write(mc *minio.Client, items []models.PictureMetadatas) error {
	var buf bytes.Buffer

	writer := parquet.NewGenericWriter[models.PictureMetadatas](&buf)
	_, err := writer.Write(items)
	if err != nil {
		log.Println("Failed to write items:", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Println("Failed to close writer:", err)
		return err
	}

	if err := utils.PutObject(mc, "database", "pictures.parquet", buf.Bytes()); err != nil {
		log.Println("Failed to put object:", err)
		return err
	}

	return nil
}

// GetAll retrieves all picture metadata from the database.
//
// This function reads the parquet file from the database bucket
// and returns all PictureMetadatas records.
//
// Parameters:
//   - mc: The MinIO client instance
//
// Returns:
//   - []models.PictureMetadatas: A slice of PictureMetadatas containing all metadata
//   - error: An error if the reading process fails
func GetAll(mc *minio.Client) ([]models.PictureMetadatas, error) {
	obj, err := utils.GetObject(mc, "database", "pictures.parquet")
	if err != nil {
		log.Println("error getting object:", err)
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("error reading object: %v", err)
	}

	reader := bytes.NewReader(data)

	items, err := parquet.Read[models.PictureMetadatas](reader, reader.Size())
	if err != nil {
		return nil, fmt.Errorf("error reading parquet: %v", err)
	}

	return items, nil
}

// Get retrieves a specific picture metadata by ID.
//
// This function is currently a stub and needs to be implemented.
//
// Parameters:
//   - id: The ID of the picture metadata to retrieve
//
// Returns:
//   - *models.PictureMetadatas: A pointer to the PictureMetadatas struct
//   - error: An error if the retrieval process fails
func Get(id uint64) (*models.PictureMetadatas, error) {
	return nil, nil
}
