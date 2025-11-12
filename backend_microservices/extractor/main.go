// Package main provides image metadata extraction functionality.
//
// This package handles the extraction of metadata from image files.
// It uses the exif package to read EXIF data from images and
// populates the PictureMetadatas struct with the extracted information.
package main

import (
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"
	"vanespen.art-microservices/common/models"
	"vanespen.art-microservices/common/utils"
)

// handleExtract processes an incoming NATS message to extract metadata from an image.
//
// This function:
// 1. Unmarshals the incoming message into a Picture struct
// 2. Retrieves the image bytes from the specified source
// 3. Extracts metadata from the image
// 4. Returns the extracted metadata or an error if any step fails
//
// Parameters:
//
//	msg: The NATS message containing the image information
//
// Returns:
//
//	models.PictureMetadatas: The extracted metadata from the image
//	error: An error if any step in the process fails
func handleExtract(msg *nats.Msg) (models.PictureMetadatas, error) {
	picture := models.Picture{}

	err := msgpack.Unmarshal(msg.Data, &picture)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return models.PictureMetadatas{}, err
	}

	log.Println("Got message:", picture)

	imageBytes, err := utils.GetImageBytes(picture)
	if err != nil {
		log.Println("Failed to get image: ", err)
		return models.PictureMetadatas{}, err
	}

	metadatas, err := extract(picture, imageBytes)
	if err != nil {
		log.Println("Failed to extract metadata: ", err)
		return models.PictureMetadatas{}, err
	}

	return metadatas, nil
}

// requestSave sends the extracted metadata to the database service.
//
// This function marshals the PictureMetadatas into msgpack format
// and sends it to the "picture.store" topic on the NATS server.
func requestSave(metadatas models.PictureMetadatas) error {
	nc := utils.NewNatsClient()
	defer nc.Close()

	payload, err := msgpack.Marshal(metadatas)
	if err != nil {
		return err
	}

	msg, err := nc.Request("picture.store", payload, 5*time.Second)
	if err != nil {
		log.Println("Failed to request store:", err)
		return err
	}

	var response models.ServiceResponse
	err = msgpack.Unmarshal(msg.Data, &response)
	if err != nil {
		log.Println("Failed to unmarshal response:", err)
		return err
	}

	if !response.Success() {
		log.Println("Failed to notify picture:", response.Msg)
		return errors.New(response.Msg)
	}

	return nil
}

// main initializes the NATS client and sets up message handlers.
//
// This function connects to the NATS server, subscribes to the
// "picture.extract" topic, and handles incoming messages to extract
// metadata from images. It also handles graceful shutdown on SIGINT
// or SIGTERM signals.
func main() {
	nc := utils.NewNatsClient()

	defer nc.Close()

	sub, err := nc.QueueSubscribe("picture.extract", "extractor_queue", func(m *nats.Msg) {
		var response models.ServiceResponse

		metadatas, extractError := handleExtract(m)
		if extractError != nil {
			response = models.ServiceResponse{
				Code: 500,
				Msg:  extractError.Error(),
			}
		} else {
			response = models.ServiceResponse{
				Code: 200,
				Msg:  "Metadata extracted successfully",
			}

			// Save the extracted metadata to the database
			requestError := requestSave(metadatas)
			if requestError != nil {
				response = models.ServiceResponse{
					Code: 500,
					Msg:  requestError.Error(),
				}
			} else {
				response = models.ServiceResponse{
					Code: 200,
					Msg:  "Metadata saved successfully",
				}
			}
		}

		message, err := msgpack.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response: ", err)
			return
		}

		err = m.Respond(message)
		if err != nil {
			log.Println("Failed to respond: ", err)
			return
		}
	})
	if err != nil {
		log.Fatal("Failed to subscribe to topic: ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	sub.Unsubscribe()
}
