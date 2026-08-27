package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultSize = 6
)

var validSlugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,50}$`)

// GenerateShortCode creates a cryptographically secure random alphanumeric string
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = defaultSize
	}

	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// IsValidSlug validates whether a custom slug contains valid characters (3-50 chars: a-z, A-Z, 0-9, _, -)
func IsValidSlug(slug string) bool {
	return validSlugRegex.MatchString(slug)
}
