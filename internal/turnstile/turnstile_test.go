package turnstile

import (
	"testing"
)

func TestVerify_NoSecretKey(t *testing.T) {
	// When no secret key is configured, verification should pass (development mode)
	t.Setenv("CLOUDFLARE_TURNSTILE_SECRET_KEY", "")

	err := Verify("any-token", "127.0.0.1")
	if err != nil {
		t.Errorf("Verify() error = %v, want nil (should skip when no secret key)", err)
	}
}

func TestVerify_NoToken(t *testing.T) {
	t.Setenv("CLOUDFLARE_TURNSTILE_SECRET_KEY", "test-secret-key")

	err := Verify("", "127.0.0.1")
	if err == nil {
		t.Error("Verify() error = nil, want error for empty token")
	}
}

func TestVerify_InvalidToken(t *testing.T) {
	t.Setenv("CLOUDFLARE_TURNSTILE_SECRET_KEY", "test-secret-key")

	// This will make a real request to Cloudflare with an invalid token
	// It should fail gracefully
	err := Verify("invalid-token", "127.0.0.1")
	if err == nil {
		t.Error("Verify() error = nil, want error for invalid token")
	}
}
