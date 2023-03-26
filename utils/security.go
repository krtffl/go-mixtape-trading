package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetEncryptionKey() ([]byte, error) {
	encodedKey := os.Getenv("ENCRYPTION_KEY")
	if encodedKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable not set")
	}
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %v", err)
	}
	return key, nil
}

func SetEncryptionKeyEnvVar() error {
	key, err := GenerateEncryptionKey()
	if err != nil {
		return errors.New("failed to generate encryption key")
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)

	if err := os.Setenv("ENCRYPTION_KEY", encodedKey); err != nil {
		return errors.New("failed to set encryption key as environment variable")
	}

	return nil
}
