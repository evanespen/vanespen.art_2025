// Package models provides data structures used across the application.
//
// This package defines the core data models that are shared between
// different microservices in the application. The models include
// picture information, metadata structures, and service response formats.
package models

// ServiceResponse represents a standard response format for service operations.
//
// This struct is used to provide consistent response formatting across
// different microservices. It includes a status code and a message
// field to indicate the result of an operation.
type ServiceResponse struct {
	// Code is the HTTP-like status code of the response
	Code int

	// Msg is a human-readable message describing the response
	Msg string
}

// Success checks if the response indicates a successful operation.
//
// A successful operation is defined as having a status code between
// 200 and 299 (inclusive).
func (s ServiceResponse) Success() bool {
	return s.Code >= 200 && s.Code < 300
}
