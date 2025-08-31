package service

import (
	"dekamond-task/internal/config"
	"testing"
)

func TestTokenService(t *testing.T) {
	cfg := &config.Config{
		Token: struct{ Secret string }{Secret: "supersecretkey"},
	}
	ts := NewTokenService(cfg)

	payload := map[string]any{"user": "john"}

	token, err := ts.GenerateToken(payload)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("expected token, got empty string")
	}

	decoded, err := ts.Decode(token)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if decoded["user"] != "john" {
		t.Fatalf("expected user=john, got %v", decoded["user"])
	}

	valid, err := ts.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if !valid {
		t.Fatal("expected token to be valid")
	}

	invalidToken := token + "x"
	valid, err = ts.Verify(invalidToken)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	if valid {
		t.Fatal("expected token to be invalid")
	}
}
