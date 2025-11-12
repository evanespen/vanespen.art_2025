// Package main provides image metadata extraction functionality.
//
// This package handles the extraction of metadata from image files.
// It uses the exif package to read EXIF data from images and
// populates the PictureMetadatas struct with the extracted information.
package main

import (
	"bytes"
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"vanespen.art-microservices/common/models"
)

type StringOrInt interface {
	string | int32 | int16
}

// computeFNumber calculates the f-number from a string representation.
//
// This function takes a string in the format "numerator/denominator"
// and returns the f-number as a float32 value.
func computeFNumber(raw string) float32 {
	parts := strings.Split(raw, "/")
	if len(parts) != 2 {
		return 0.0
	}
	numerator, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return 0.0
	}
	denominator, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return 0.0
	}
	if denominator == 0 {
		return 0.0
	}
	return float32(numerator / denominator)
}

// computeFocalLength calculates the focal length from a string representation.
//
// This function is a wrapper around computeFNumber for focal length values.
func computeFocalLength(raw string) float32 {
	return computeFNumber(raw)
}

// computeMode converts a raw exposure program value to a human-readable mode.
//
// This function maps numeric exposure program values to descriptive strings.
func computeMode(raw string) string {
	switch raw {
	case "0":
		return "Manual"
	case "1":
		return "Program"
	case "2":
		return "Aperture Priority"
	case "3":
		return "Shutter Priority"
	case "4":
		return "Creative"
	case "5":
		return "Action"
	case "6":
		return "Portrait"
	case "7":
		return "Landscape"
	default:
		return "Unknown"
	}
}

// getExifValue retrieves a value from EXIF data based on field name and type.
//
// This function extracts a value from EXIF data and converts it to the
// specified type (string or int).
func getExifValue[T string | int](fieldName exif.FieldName, exifData *exif.Exif) T {
	value, err := exifData.Get(fieldName)
	if err != nil {
		fmt.Println("Failed to get exif value:", err)
		var zero T
		return zero
	}

	switch any(*new(T)).(type) {
	case string:
		return any(strings.ReplaceAll(value.String(), "\"", "")).(T)
	case int:
		intValue, _ := strconv.Atoi(value.String())
		return any(intValue).(T)
	default:
		var zero T
		return zero
	}
}

// extract processes an image and extracts its metadata.
//
// This function takes a Picture and its byte data, extracts EXIF metadata,
// and returns a populated PictureMetadatas struct.
func extract(picture models.Picture, data []byte) (models.PictureMetadatas, error) {

	exifData, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return models.PictureMetadatas{}, err
	}

	imageDatetime, _ := exifData.DateTime()

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return models.PictureMetadatas{}, err
	}

	height := img.Bounds().Dy()
	width := img.Bounds().Dx()

	return models.PictureMetadatas{
		Id:             picture.Key,
		Ext:            picture.Ext,
		Timestamp:      int64(imageDatetime.Unix()),
		Camera:         getExifValue[string](exif.FieldName("Model"), exifData),
		Mode:           computeMode(getExifValue[string](exif.FieldName("ExposureProgram"), exifData)),
		Aperture:       computeFNumber(getExifValue[string](exif.FieldName("FNumber"), exifData)),
		Iso:            int32(getExifValue[int](exif.FieldName("ISOSpeedRatings"), exifData)),
		Speed:          getExifValue[string](exif.FieldName("ExposureTime"), exifData),
		FocalLength:    computeFocalLength(getExifValue[string](exif.FieldName("FocalLength"), exifData)),
		Lens:           getExifValue[string](exif.FieldName("LensModel"), exifData),
		Flash:          false,
		Landscape:      width > height,
		Panoramic:      width > height && float32(width)/float32(height) > 1.5,
		Width:          int32(width),
		Height:         int32(height),
		Favourite:      false,
		TriggerWarning: false,
		Description:    "",
	}, nil
}
