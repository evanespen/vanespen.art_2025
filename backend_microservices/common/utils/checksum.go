// Package utils provides utility functions for various operations.
//
// This package contains helper functions that are used across
// different microservices. It includes utilities for MinIO
// object storage and NATS messaging.
package utils

import "github.com/cespare/xxhash"

// Hash computes a 64-bit hash of the input data using the xxHash algorithm.
//
// This function is a wrapper around xxhash.Sum64 that provides a consistent
// interface for hashing byte slices. The xxHash algorithm is known for its
// speed and good hash distribution properties.
//
// Parameters:
//   - data: the byte slice to be hashed
//
// Returns:
//   - uint64: the 64-bit hash value of the input data
func Hash(data []byte) uint64 {
	return xxhash.Sum64(data)
}
