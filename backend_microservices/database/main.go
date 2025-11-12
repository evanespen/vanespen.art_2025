// Package main provides database operations for the application.
//
// This package handles all database-related functionality, including
// connecting to the database, performing CRUD operations, and managing
// data storage. It uses MinIO for object storage and parquet format
// for efficient data serialization.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"
	"vanespen.art-microservices/common/models"
	"vanespen.art-microservices/common/utils"
)

// handleStore processes incoming store requests from NATS.
//
// This function handles messages from the "picture.store" topic,
// unmarshals the PictureMetadatas from the message, and writes
// them to the database.
//
// Parameters:
//   - mc: The MinIO client instance
//   - msg: The NATS message containing the picture metadata
//
// Returns:
//   - error: An error if the processing or writing fails
func handleStore(mc *minio.Client, msg *nats.Msg) error {
	metadata := models.PictureMetadatas{}
	if err := msgpack.Unmarshal(msg.Data, &metadata); err != nil {
		return err
	}

	log.Println("Got message:", metadata)

	// err := Write(mc, []models.PictureMetadatas{metadata})

	var response models.ServiceResponse
	if insertError := Insert(mc, metadata); insertError != nil {
		log.Println("Failed to insert", metadata, ":", insertError)
		response.Code = 500
		response.Msg = "Failed to insert metadata"
	} else {
		response.Code = 200
		response.Msg = "Metadata inserted successfully"
	}

	message, err := msgpack.Marshal(response)
	if err != nil {
		log.Println("Got marshal error:", err)
		return err
	}

	if err := msg.Respond(message); err != nil {
		log.Println("Got respond error:", err)
		return err
	}

	return nil
}

// handleGetAll processes incoming get_all requests from NATS.
//
// This function handles messages from the "picture.get_all" topic,
// retrieves all PictureMetadatas from the database, and responds
// with the serialized data.
//
// Parameters:
//   - mc: The MinIO client instance
//   - msg: The NATS message triggering the retrieval
//
// Returns:
//   - error: An error if the retrieval or response fails
func handleGetAll(mc *minio.Client, msg *nats.Msg) error {
	items, err := GetAll(mc)
	if err != nil {
		return err
	}

	response, err := msgpack.Marshal(items)
	if err != nil {
		log.Println("Got marshal error:", err)
		return err
	}

	err = msg.Respond(response)
	if err != nil {
		log.Println("Got respond error:", err)
		return err
	}

	return nil
}

// main initializes the NATS client and sets up message handlers.
//
// This function connects to the NATS server, subscribes to the
// "picture.store" and "picture.get_all" topics, and handles
// graceful shutdown on SIGINT or SIGTERM signals.
func main() {
	nc := utils.NewNatsClient()
	defer nc.Close()

	mc := utils.NewMinioClient()

	storeSub, storeSubErr := nc.QueueSubscribe("picture.store", "database_queue_store", func(m *nats.Msg) {
		handleStore(mc, m)
	})
	if storeSubErr != nil {
		log.Fatal("Failed to subscribe to topic: ", storeSubErr)
	}

	getSub, getSubErr := nc.QueueSubscribe("picture.get_all", "database_queue_get", func(m *nats.Msg) {
		handleGetAll(mc, m)
	})
	if getSubErr != nil {
		log.Fatal("Failed to subscribe to topic: ", getSubErr)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	storeSub.Unsubscribe()
	getSub.Unsubscribe()
}
