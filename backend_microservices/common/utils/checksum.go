package utils

import "github.com/cespare/xxhash"

func Hash(data []byte) uint64 {
	return xxhash.Sum64(data)
}
