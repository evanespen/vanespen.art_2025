// Package main provides image resizing functionality.

// This package handles the resizing of images to different sizes.
// It listens for messages on the NATS "picture.resize" topic,
// processes the images, creates resized versions (half and thumb),
// and stores them in MinIO buckets.
package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"

	"github.com/disintegration/imaging"
	"vanespen.art-microservices/common/models"
	"vanespen.art-microservices/common/utils"
)

// resizeImage resizes an image by a specified factor.
//
// This function takes an image and resizes it to the specified factor
// of its original dimensions. It supports both JPEG and PNG formats.
func resizeImage(picture models.Picture, img image.Image, factor float64) ([]byte, error) {
	// Implement image resizing logic here
	log.Println("Resizing image:", picture.Key)

	resizedImg := imaging.Resize(img, int(float64(img.Bounds().Dx())*factor), int(float64(img.Bounds().Dy())*factor), imaging.Lanczos)

	buf := bytes.NewBuffer(nil)

	switch picture.Ext {
	case "jpg", "jpeg":
		jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(buf, resizedImg)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", picture.Ext)
	}

	return buf.Bytes(), nil
}

// saveImage saves a resized image to MinIO storage.
//
// This function saves the resized image to the appropriate MinIO bucket
// (either "half" or "thumb" based on the kind parameter).
func saveImage(picture models.Picture, imageBytes []byte, kind string) error {
	// Implement image saving logic here
	log.Println("Saving resized image:", picture.Key)

	var bucketName string
	if kind != "half" && kind != "thumb" {
		return fmt.Errorf("Unknown kind: %s", kind)
	} else {
		bucketName = kind
	}

	minioClient := utils.NewMinioClient()

	// Check if the bucket exists, if not, create it
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("Error checking bucket existence: %s", err)

	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("Error creating bucket: %s", err)
		}
	}

	objectName := fmt.Sprintf("%d.%s", picture.Key, picture.Ext)

	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		bytes.NewReader(imageBytes),
		int64(len(imageBytes)),
		minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("Failed to save %s image: %s", kind, err)
	}

	return nil
}

// messageHandler processes incoming resize requests from NATS.
//
// This function handles messages from the "picture.resize" topic,
// retrieves the image from MinIO, resizes it to half and thumb sizes,
// and saves the resized images back to MinIO.
func messageHandler(msg *nats.Msg) error {
	picture := models.Picture{}

	err := msgpack.Unmarshal(msg.Data, &picture)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return err
	}

	img, err := utils.GetImage(picture)
	if err != nil {
		log.Println("Failed to get image:", err)
		return err
	}

	halfImage, err := resizeImage(picture, img, 0.5)
	if err != nil {
		log.Println("Failed to create half image:", err)
		return err
	}

	err = saveImage(picture, halfImage, "half")
	if err != nil {
		log.Println("Failed to save half image:", err)
		return err
	}

	thumbImage, err := resizeImage(picture, img, 0.1)
	if err != nil {
		log.Println("Failed to create half image:", err)
		return err
	}

	err = saveImage(picture, thumbImage, "thumb")
	if err != nil {
		log.Println("Failed to save thumb image:", err)
		return err
	}

	log.Println("Image processing completed")
	return nil
	// Notify
}

// main initializes the NATS client and sets up message handlers.
//
// This function connects to the NATS server, subscribes to the
// "picture.resize" topic, and handles incoming messages to resize
// images. It also handles graceful shutdown on SIGINT or SIGTERM signals.
func main() {
	// connect to the nats server
	nc := utils.NewNatsClient()
	defer nc.Close()

	sub, err := nc.Subscribe("picture.resize", func(m *nats.Msg) {
		resizeError := messageHandler(m)
		var response models.ServiceResponse

		if resizeError != nil {
			log.Println("Failed to resize image:", resizeError)
			response = models.ServiceResponse{
				Code: 500,
				Msg:  resizeError.Error(),
			}
		} else {
			response = models.ServiceResponse{
				Code: 200,
				Msg:  "Image resized successfully",
			}
		}

		responsesBytes, err := msgpack.Marshal(response)
		if err != nil {
			log.Fatal("Failed to marshal response:", err)
		}

		if err := m.Respond(responsesBytes); err != nil {
			log.Fatal("Failed to respond:", err)
		}
	})
	if err != nil {
		log.Fatal("Failed to subscribe: ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	sub.Unsubscribe()
}
