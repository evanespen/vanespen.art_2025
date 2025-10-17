package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateSHA256Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file : %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("cannot compute checksum : %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
