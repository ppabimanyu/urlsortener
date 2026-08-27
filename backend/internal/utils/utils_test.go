package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPass123!"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hashed) {
		t.Fatalf("Password verification failed for correct password")
	}

	if CheckPasswordHash("wrongPassword", hashed) {
		t.Fatalf("Password verification succeeded for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	secret := "test_secret_key_12345"
	userID := uuid.New()
	email := "test@example.com"
	name := "Test User"

	tokenStr, err := GenerateToken(userID, email, name, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(tokenStr, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID || claims.Email != email || claims.Name != name {
		t.Fatalf("Claims mismatch: got %+v", claims)
	}

	// Validate with wrong secret
	_, err = ValidateToken(tokenStr, "wrong_secret")
	if err == nil {
		t.Fatalf("Token validation should have failed with wrong secret")
	}
}

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode(6)
	if err != nil {
		t.Fatalf("Failed to generate short code: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("Expected length 6, got %d", len(code))
	}

	if !IsValidSlug(code) {
		t.Fatalf("Generated code %s failed IsValidSlug", code)
	}
}

func TestIsValidSlug(t *testing.T) {
	validSlugs := []string{"promo-2026", "my_link", "abc123", "Sale-Today_50"}
	for _, slug := range validSlugs {
		if !IsValidSlug(slug) {
			t.Errorf("Slug %s should be valid", slug)
		}
	}

	invalidSlugs := []string{"ab", "invalid slug with spaces", "special!@#", "a/b"}
	for _, slug := range invalidSlugs {
		if IsValidSlug(slug) {
			t.Errorf("Slug %s should be invalid", slug)
		}
	}
}

func TestParseUserAgent(t *testing.T) {
	uaChromeMac := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	device, os, browser := ParseUserAgent(uaChromeMac)

	if device != "Desktop" {
		t.Errorf("Expected Desktop device, got %s", device)
	}
	if os != "macOS" {
		t.Errorf("Expected macOS, got %s", os)
	}
	if browser != "Chrome" {
		t.Errorf("Expected Chrome, got %s", browser)
	}
}
