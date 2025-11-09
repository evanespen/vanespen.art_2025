package main

import (
	"context"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack"
)

type Picture struct {
	Key        uuid.UUID
	Ext        string
	BytesCount int64
}

type ServiceResponse struct {
	Success bool
	Msg     string
}

func getMinioClient() *minio.Client {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Error creating Minio client:", err)
		return nil
	}
	return minioClient
}

func getImage(picture Picture) ([]byte, error) {
	log.Println("Retrieving image:", picture.Key)

	minioClient := getMinioClient()

	bucketName := "full"
	objectName := fmt.Sprintf("%s.%s", picture.Key, picture.Ext)
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

	// // load the bytes to an image
	// img, _, err := image.Decode(bytes.NewReader(imageBytes))
	// if err != nil {
	// 	log.Fatal("Error decoding image:", err)
	// 	return nil, err
	// }

	return imageBytes, nil
}

func handleExtract(msg *nats.Msg) error {
	picture := Picture{}

	err := msgpack.Unmarshal(msg.Data, &picture)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return err
	}

	log.Println("Got message:", picture)

	img, err := getImage(picture)
	if err != nil {
		log.Println("Failed to get image: ", err)
		return err
	}

	extract(picture.Key.String(), img)

	return nil
}

func main() {
	// connect to the nats server
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to NATS server")

	defer nc.Close()

	sub, err := nc.Subscribe("picture.extract", func(m *nats.Msg) {
		var response ServiceResponse

		extractError := handleExtract(m)
		if extractError != nil {
			response = ServiceResponse{
				Success: false,
				Msg:     extractError.Error(),
			}
		} else {
			response = ServiceResponse{
				Success: true,
				Msg:     "metadatas:",
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	sub.Unsubscribe()
}
