package main

import (
	"bytes"
	"log"

	"github.com/rwcarlsen/goexif/exif"
)

type PictureMetadatas struct {
	UUID           string `json:"uuid" parquet:"name=uuid, type=FIXED_LEN_BYTE_ARRAY, length=37, convertedtype=UTF8, encoding=DELTA_BYTE_ARRAY"`
	Ext            string `json:"ext" parquet:"name=ext, type=BYTE_ARRAY, convertedtype=UTF8, encoding=DELTA_BYTE_ARRAY"`
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
	Panoramic      bool   `json:"panoramic" parquet:"name=panoramic, type=BOOLEAN, encoding=PLAIN"`
	Width          int16  `json:"width" parquet:"name=width, type=INT32, encoding=DELTA_BINARY_PACKED"`
	Height         int16  `json:"height" parquet:"name=height, type=INT32, encoding=DELTA_BINARY_PACKED"`
	Favourite      bool   `json:"favourite" parquet:"name=favourite, type=BOOLEAN, encoding=PLAIN"`
	TriggerWarning bool   `json:"trigger_warning" parquet:"name=trigger_warning, type=BOOLEAN, encoding=PLAIN"`
	Description    string `json:"description" parquet:"name=description, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
}

func extract(key string, data []byte) (PictureMetadatas, error) {

	exifData, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return PictureMetadatas{}, err
	}

	log.Println(exifData)

	return PictureMetadatas{}, nil

	// var pictureMetadatas PictureMetadatas
	// pictureMetadatas.UUID = "uuid"         // Placeholder, replace with actual UUID extraction
	// pictureMetadatas.Ext = "ext"           // Placeholder, replace with actual extension extraction
	// pictureMetadatas.Checksum = "checksum" // Placeholder, replace with actual checksum extraction
	// pictureMetadatas.Timestamp = int(exifData.DateTime.Unix())
	// pictureMetadatas.Camera = exifData.Model
	// pictureMetadatas.Mode = "mode" // Placeholder, replace with actual mode extraction
	// pictureMetadatas.Aperture = fmt.Sprintf("f/%f", exifData.ApertureValue)
	// pictureMetadatas.Iso = int16(exifData.ISOSpeedRatings[0])
	// pictureMetadatas.Speed = exifData.ShutterSpeedValue.String()
	// pictureMetadatas.FocalLength = exifData.FocalLength.String()
	// pictureMetadatas.Lens = "lens" // Placeholder, replace with actual lens extraction
	// pictureMetadatas.Flash = exifData.Flash.Fired
	// pictureMetadatas.Landscape = exifData.ImageWidth > exifData.ImageHeight
	// pictureMetadatas.Panoramic = pictureMetadatas.Landscape && float32(exifData.ImageWidth)/float32(exifData.ImageHeight) > 1.5
	// pictureMetadatas.Width = int16(exifData.ImageWidth)
	// pictureMetadatas.Height = int16(exifData.ImageHeight)
	// pictureMetadatas.Favourite = false
	// pictureMetadatas.TriggerWarning = false
	// pictureMetadatas.Description = "description" // Placeholder, replace with actual description extraction

	// return pictureMetadatas, nil

}

// func NewPicture(imagePath string, pictureUUID string) (PictureMetadatas, error) {
// 	et, err := exiftool.NewExiftool()
// 	if err != nil {
// 		fmt.Printf("Error when intializing: %v\n", err)
// 		return PictureMetadatas{}, errors.New(fmt.Sprintf("Error when intializing: %v\n", err))
// 	}
// 	defer et.Close()

// 	fileInfos := et.ExtractMetadata(imagePath)
// 	extension := path.Ext(imagePath)

// 	datetime, err := time.Parse("2006:01:02 15:04:05", fileInfos[0].Fields["DateTimeOriginal"].(string))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	timestamp := datetime.Unix()

// 	// checksum, _ := utils.CalculateSHA256Checksum(imagePath)
// 	checksum := pictureUUID

// 	width := int16(fileInfos[0].Fields["ImageWidth"].(float64))
// 	height := int16(fileInfos[0].Fields["ImageHeight"].(float64))
// 	landscape := width > height

// 	var panoramic bool
// 	if landscape && float32(width)/float32(height) > 1.5 {
// 		panoramic = true
// 	} else {
// 		panoramic = false
// 	}

// 	return PictureMetadatas{
// 		UUID:           pictureUUID,
// 		Ext:            extension,
// 		Checksum:       checksum,
// 		Camera:         fileInfos[0].Fields["Model"].(string),
// 		Timestamp:      int(timestamp),
// 		Mode:           fileInfos[0].Fields["ExposureProgram"].(string),
// 		Aperture:       fmt.Sprintf("f/%f", math.Round(fileInfos[0].Fields["FNumber"].(float64))),
// 		Iso:            int16(math.Round(fileInfos[0].Fields["ISO"].(float64))),
// 		Speed:          fileInfos[0].Fields["ShutterSpeed"].(string),
// 		FocalLength:    fileInfos[0].Fields["FocalLength"].(string),
// 		Lens:           fileInfos[0].Fields["LensID"].(string),
// 		Flash:          fileInfos[0].Fields["Flash"].(string) == "Off, Did not fire.",
// 		Landscape:      landscape,
// 		Panoramic:      panoramic,
// 		Width:          width,
// 		Height:         height,
// 		Favourite:      false,
// 		TriggerWarning: false,
// 		Description:    "",
// 	}, nil
// }
