// Package utils provides utility functions for various operations.
//
// This package contains helper functions that are used across
// different microservices. It includes utilities for MinIO
// object storage and NATS messaging.
package utils

import (
	"log"

	"github.com/nats-io/nats.go"
)

// NewNatsClient creates a new NATS client instance.
//
// This function initializes and returns a NATS client configured
// to connect to the local NATS server at the default address.
//
// Returns:
// - *nats.Conn: A pointer to the initialized NATS client.
func NewNatsClient() *nats.Conn {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to NATS server")

	return nc
}
