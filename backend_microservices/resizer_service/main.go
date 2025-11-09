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

	"github.com/disintegration/imaging"
)

// Picture represents a picture stored in MinIO with its metadata
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

func getImage(picture Picture) (image.Image, error) {
	log.Println("Got message:", picture)
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

	// load the bytes to an image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal("Error decoding image:", err)
		return nil, err
	}

	return img, nil
}

func resizeImage(picture Picture, img image.Image, factor float64) ([]byte, error) {
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

func saveImage(picture Picture, imageBytes []byte, kind string) error {
	// Implement image saving logic here
	log.Println("Saving resized image:", picture.Key)

	var bucketName string
	if kind != "half" && kind != "thumb" {
		return fmt.Errorf("Unknown kind:", kind)
	} else {
		bucketName = kind
	}

	minioClient := getMinioClient()

	// Check if the bucket exists, if not, create it
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("Error checking bucket existence:", err)

	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("Error creating bucket:", err)
		}
	}

	objectName := fmt.Sprintf("%s.%s", picture.Key, picture.Ext)

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

func messageHandler(msg *nats.Msg) error {
	picture := Picture{}

	err := msgpack.Unmarshal(msg.Data, &picture)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return err
	}

	img, err := getImage(picture)
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

func main() {
	// connect to the nats server
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to NATS server")

	defer nc.Close()

	sub, err := nc.Subscribe("picture.resize", func(m *nats.Msg) {
		resizeError := messageHandler(m)
		var response ServiceResponse

		if resizeError != nil {
			log.Println("Failed to resize image:", resizeError)
			response = ServiceResponse{
				Success: false,
				Msg:     resizeError.Error(),
			}
		} else {
			response = ServiceResponse{
				Success: true,
				Msg:     "Image resized successfully",
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
