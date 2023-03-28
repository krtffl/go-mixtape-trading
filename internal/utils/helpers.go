package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetEncryptionKey() ([]byte, error) {
	encodedKey := os.Getenv("ENCRYPTION_KEY")
	if encodedKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY does not exist on env variables")
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %v", err)
	}
	return key, nil
}

func SetEncryptionKeyEnvVar() error {
	log.Print("generating cookie encryption key")

	key, err := generateEncryptionKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %v", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	if err := os.Setenv("ENCRYPTION_KEY", encodedKey); err != nil {
		return fmt.Errorf("failed to set encryption key as environment variable: %v", err)
	}

	log.Print("encryption key stored as environment variable")

	return nil
}

func SetCookie(w http.ResponseWriter, name string, value string) error {
	log.Printf("setting cookie %s", name)

	key, err := GetEncryptionKey()
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(value))
	stream.XORKeyStream(ciphertext, []byte(value))

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	encodedIV := base64.StdEncoding.EncodeToString(iv)

	cookie := &http.Cookie{
		Name:  name,
		Value: encodedCiphertext + "|" + encodedIV,
	}
	http.SetCookie(w, cookie)

	log.Printf("cookie %s set", name)

	return nil
}

func GetCookie(r *http.Request, name string) (string, error) {
	log.Printf("retrieving cookie %s", name)

	key, err := GetEncryptionKey()
	if err != nil {
		return "", err
	}

	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return "", errors.New("invalid cookie value")
	}
	encodedCiphertext := parts[0]
	encodedIV := parts[1]

	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(encodedIV)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	log.Printf("cookie %s successfully retrieved", name)

	return string(plaintext), nil
}

func generateEncryptionKey() ([]byte, error) {
	log.Print("generating encryption key")

	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		return nil, err
	}

	log.Printf("encryption key successfully generated")

	return key, nil
}
