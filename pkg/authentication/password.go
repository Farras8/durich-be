package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	MaxPasswordLength = 72
	DefaultCost       = 12
)

var (
	ErrPasswordTooLong = errors.New("password is too long")
)

func HashPassword(password string) (string, error) {
	if err := validatePassword(password); err != nil {
		return "", err
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("error generating salt: %w", err)
	}

	saltedPassword := []byte(password + string(salt))

	hash, err := bcrypt.GenerateFromPassword(saltedPassword, DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", encodedSalt, encodedHash), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	if err := validatePassword(password); err != nil {
		return false, err
	}

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("error decoding salt: %w", err)
	}

	hash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %w", err)
	}

	saltedPassword := []byte(password + string(salt))

	err = bcrypt.CompareHashAndPassword(hash, saltedPassword)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error comparing password and hash: %w", err)
	}

	return true, nil
}

func validatePassword(password string) error {
	if len(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}
	return nil
}