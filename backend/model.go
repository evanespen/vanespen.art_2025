package main

import (
	"errors"
	"fmt"
	"github.com/barasher/go-exiftool"
	"github.com/google/uuid"
	"math"
	"time"
)

type Album struct {
	UUID        string   `json:"uuid" parquet:"name=uuid, type=FIXED_LEN_BYTE_ARRAY, length=37, convertedtype=UTF8, encoding=DELTA_BYTE_ARRAY"`
	Title       string   `json:"title" parquet:"name=title, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Description string   `json:"description" parquet:"name=description, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Pictures    []string `json:"pictures" parquet:"name=pictures, type=MAP convertedtype=LIST, valuetype=FIXED_LEN_BYTE_ARRAY, valueconvertedtype=UTF8, valuelength=37"`
}

type Picture struct {
	UUID           string `json:"uuid" parquet:"name=uuid, type=FIXED_LEN_BYTE_ARRAY, length=37, convertedtype=UTF8, encoding=DELTA_BYTE_ARRAY"`
	Checksum       string `json:"checksum" parquet:"name=checksum, type=FIXED_LEN_BYTE_ARRAY, length=65, encoding=DELTA_BYTE_ARRAY"`
	Timestamp      int    `json:"timestamp" parquet:"name=timestamp, type=INT64, encoding=DELTA_BINARY_PACKED"`
	Camera         string `json:"camera" parquet:"name=camera, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Mode           string `json:"mode" parquet:"name=mode, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Aperture       string `json:"aperture" parquet:"name=aperture, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Iso            int16  `json:"iso" parquet:"name=iso, type=INT32, encoding=DELTA_BINARY_PACKED"`
	Speed          string `json:"speed" parquet:"name=speed, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	FocalLength    string `json:"focal_length" parquet:"name=focal_length, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Lens           string `json:"lens" parquet:"name=lens, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Flash          bool   `json:"flash" parquet:"name=flash, type=BOOLEAN, encoding=PLAIN"`
	Landscape      bool   `json:"landscape" parquet:"name=landscape, type=BOOLEAN, encoding=PLAIN"`
	Width          int16  `json:"width" parquet:"name=width, type=INT32, encoding=DELTA_BINARY_PACKED"`
	Height         int16  `json:"height" parquet:"name=height, type=INT32, encoding=DELTA_BINARY_PACKED"`
	Favourite      bool   `json:"favourite" parquet:"name=favourite, type=BOOLEAN, encoding=PLAIN"`
	TriggerWarning bool   `json:"trigger_warning" parquet:"name=trigger_warning, type=BOOLEAN, encoding=PLAIN"`
	Description    string `json:"description" parquet:"name=description, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
}

func NewPicture(imagePath string, pictureUUID string) (Picture, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return Picture{}, errors.New(fmt.Sprintf("Error when intializing: %v\n", err))
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(imagePath)

	datetime, err := time.Parse("2006:01:02 15:04:05", fileInfos[0].Fields["DateTimeOriginal"].(string))
	if err != nil {
		fmt.Println(err)
	}
	timestamp := datetime.Unix()

	checksum, _ := CalculateSHA256Checksum(imagePath)

	//for _, fileInfo := range fileInfos {
	//	if fileInfo.Err != nil {
	//		fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
	//		continue
	//	}
	//
	//	for k, v := range fileInfo.Fields {
	//		fmt.Printf("[%v] %v\n", k, v)
	//	}
	//}

	return Picture{
		UUID:           pictureUUID,
		Checksum:       checksum,
		Camera:         fileInfos[0].Fields["Model"].(string),
		Timestamp:      int(timestamp),
		Mode:           fileInfos[0].Fields["ExposureProgram"].(string),
		Aperture:       fmt.Sprintf("f/%f", math.Round(fileInfos[0].Fields["FNumber"].(float64))),
		Iso:            int16(math.Round(fileInfos[0].Fields["ISO"].(float64))),
		Speed:          fileInfos[0].Fields["ShutterSpeed"].(string),
		FocalLength:    fileInfos[0].Fields["FocalLength"].(string),
		Lens:           fileInfos[0].Fields["LensID"].(string),
		Flash:          fileInfos[0].Fields["Flash"].(string) == "Off, Did not fire.",
		Landscape:      fileInfos[0].Fields["ImageWidth"].(float64) > fileInfos[0].Fields["ImageHeight"].(float64),
		Width:          int16(fileInfos[0].Fields["ImageWidth"].(float64)),
		Height:         int16(fileInfos[0].Fields["ImageHeight"].(float64)),
		Favourite:      false,
		TriggerWarning: false,
		Description:    "",
	}, nil
}

func NewAlbum() *Album {
	return &Album{
		UUID:     uuid.New().String(),
		Pictures: make([]string, 0), // Initialise un slice de chaînes vide, non-nil.
	}
}
