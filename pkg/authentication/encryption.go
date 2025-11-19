package authentication

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

var byteKey [32]byte
var onceEncrypt = &sync.Once{}

func SetupKey(key string) {
	onceEncrypt.Do(func() {
		byteKey = sha256.Sum256([]byte(key))
	})
}

func Encrypt(plaintext string) (string, error) {
	var zeroArray [32]byte
	if bytes.Equal(byteKey[:], zeroArray[:]) {
		return "", errors.New("byteKey is nil")
	}

	byteText := []byte(plaintext)
	block, err := aes.NewCipher(byteKey[:])
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(byteText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], byteText)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encodedCiphertext string) (*string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(byteKey[:])
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	plaintext := string(ciphertext)
	return &plaintext, nil
}
