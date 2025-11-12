// Package models provides data structures used across the application.
//
// This package defines the core data models that are shared between
// different microservices in the application. The models include
// picture information, metadata structures, and service response formats.
package models

// Picture represents a basic picture entity with essential information.
//
// This struct contains the minimum required fields to identify and
// describe a picture in the system. It's used for initial uploads
// and basic operations where detailed metadata isn't needed.
type Picture struct {
	// Key is a unique identifier for the picture
	Key uint64

	// Ext is the file extension of the picture (e.g., "jpg", "png")
	Ext string

	// BytesCount is the size of the picture file in bytes
	BytesCount int64
}
