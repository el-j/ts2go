package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestService_GenerateAndValidateToken(t *testing.T) {
	secret := "test-secret"
	svc := NewService(secret, time.Hour, 24*time.Hour)

	userID := uuid.New()
	email := "test@example.com"
	username := "testuser"

	// Test GenerateToken
	token, err := svc.GenerateToken(userID, email, username)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if len(strings.Split(token, ".")) != 3 {
		t.Errorf("Expected valid JWT token format")
	}

	// Test ValidateToken
	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}

	// Test invalid token
	_, err = svc.ValidateToken("invalid.token.string")
	if err == nil {
		t.Errorf("Expected error for invalid token")
	}

	// Test tampered token
	parts := strings.Split(token, ".")
	tampered := parts[0] + "." + parts[1] + "a." + parts[2]
	_, err = svc.ValidateToken(tampered)
	if err == nil {
		t.Errorf("Expected error for tampered token")
	}
}

func TestService_TokenExpiry(t *testing.T) {
	secret := "test-secret"
	// Create service with negative expiry so tokens expire immediately
	svc := NewService(secret, -1*time.Hour, -24*time.Hour)

	userID := uuid.New()

	token, err := svc.GenerateToken(userID, "test@example.com", "testuser")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Errorf("Expected error for expired token")
	}
}

func TestService_PasswordHashing(t *testing.T) {
	svc := NewService("secret", time.Hour, time.Hour)
	password := "my_secure_password"

	hash, err := svc.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == password {
		t.Errorf("Hash should not equal plain text password")
	}

	err = svc.VerifyPassword(password, hash)
	if err != nil {
		t.Errorf("VerifyPassword failed for correct password: %v", err)
	}

	err = svc.VerifyPassword("wrong_password", hash)
	if err == nil {
		t.Errorf("VerifyPassword should fail for incorrect password")
	}
}

func TestService_APIKeys(t *testing.T) {
	svc := NewService("secret", time.Hour, time.Hour)

	key, err := svc.GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey failed: %v", err)
	}
	if !strings.HasPrefix(key, "ts2go_") {
		t.Errorf("Expected API key to start with ts2go_, got %s", key)
	}

	hash, err := svc.HashAPIKey(key)
	if err != nil {
		t.Fatalf("HashAPIKey failed: %v", err)
	}

	err = svc.VerifyAPIKey(key, hash)
	if err != nil {
		t.Errorf("VerifyAPIKey failed for correct key: %v", err)
	}

	err = svc.VerifyAPIKey("ts2go_invalid", hash)
	if err == nil {
		t.Errorf("VerifyAPIKey should fail for incorrect key")
	}
}
