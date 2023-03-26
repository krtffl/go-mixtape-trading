package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
)

func SetCookie(w http.ResponseWriter, name string, value string) error {
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

	return nil
}

func GetCookie(r *http.Request, name string) (string, error) {
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

	return string(plaintext), nil
}
