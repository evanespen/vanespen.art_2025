// Package models provides data structures used across the application.
//
// This package defines the core data models that are shared between
// different microservices in the application. The models include
// picture information, metadata structures, and service response formats.
package models

// PictureMetadatas represents metadata information for a picture.
// It includes technical details about the image capture and its properties.
type PictureMetadatas struct {
	Id             uint64  `json:"id" parquet:"id"`
	Ext            string  `json:"ext" parquet:"ext"`
	Timestamp      int64   `json:"timestamp" parquet:"timestamp,timestamp"`
	Camera         string  `json:"camera" parquet:"camera"`
	Mode           string  `json:"mode" parquet:"mode"`
	Aperture       float32 `json:"aperture" parquet:"aperture"`
	Iso            int32   `json:"iso" parquet:"iso"`
	Speed          string  `json:"speed" parquet:"speed"`
	FocalLength    float32 `json:"focal_length" parquet:"focal_length"`
	Lens           string  `json:"lens" parquet:"lens"`
	Flash          bool    `json:"flash" parquet:"flash"`
	Landscape      bool    `json:"landscape" parquet:"landscape"`
	Panoramic      bool    `json:"panoramic" parquet:"panoramic"`
	Width          int32   `json:"width" parquet:"width"`
	Height         int32   `json:"height" parquet:"height"`
	Favourite      bool    `json:"favourite" parquet:"favourite"`
	TriggerWarning bool    `json:"trigger_warning" parquet:"trigger_warning"`
	Description    string  `json:"description" parquet:"description"`
}
