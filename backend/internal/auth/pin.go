package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 3
	argon2Memory  = 65536
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16
)

var (
	ErrInvalidHash = errors.New("invalid hash format")
	ErrInvalidPIN  = errors.New("pin must be exactly 6 numeric digits")
)

const pinLen = 6

func ValidatePIN(pin string) error {
	if len(pin) != pinLen {
		return ErrInvalidPIN
	}
	for _, r := range pin {
		if !unicode.IsDigit(r) {
			return ErrInvalidPIN
		}
	}
	return nil
}

func HashPin(pin string, time, memory uint32, threads uint8) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(pin), salt, time, memory, threads, argon2KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, b64Salt, b64Hash), nil
}

func VerifyPin(pin, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, ErrInvalidHash
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(pin), salt, time, memory, threads, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
