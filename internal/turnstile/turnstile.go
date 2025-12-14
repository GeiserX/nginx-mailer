package turnstile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const verifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type VerifyResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
	ChallengeTS string   `json:"challenge_ts,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
}

func Verify(token string, remoteIP string) error {
	secretKey := os.Getenv("CLOUDFLARE_TURNSTILE_SECRET_KEY")

	// If no secret key configured, skip verification (for development)
	if secretKey == "" {
		return nil
	}

	// If no token provided, fail
	if token == "" {
		return fmt.Errorf("no turnstile token provided")
	}

	// Prepare request
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)
	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	// Make request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(verifyURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to verify turnstile: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse turnstile response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("turnstile verification failed: %v", result.ErrorCodes)
	}

	return nil
}
