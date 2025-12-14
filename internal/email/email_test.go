package email

import (
	"strings"
	"testing"
)

func TestBuildEmailBody(t *testing.T) {
	data := ContactEmail{
		Name:     "Test User",
		Phone:    "+34 666 777 888",
		Email:    "test@example.com",
		Location: "Madrid",
		Message:  "This is a test message",
	}

	body := buildEmailBody(data)

	// Check that all fields are present
	if !strings.Contains(body, "Test User") {
		t.Error("Email body should contain name")
	}
	if !strings.Contains(body, "+34 666 777 888") {
		t.Error("Email body should contain phone")
	}
	if !strings.Contains(body, "test@example.com") {
		t.Error("Email body should contain email")
	}
	if !strings.Contains(body, "Madrid") {
		t.Error("Email body should contain location")
	}
	if !strings.Contains(body, "This is a test message") {
		t.Error("Email body should contain message")
	}

	// Check HTML structure
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Email body should be HTML")
	}
	if !strings.Contains(body, "</html>") {
		t.Error("Email body should have closing HTML tag")
	}
}

func TestBuildEmailBody_OptionalFields(t *testing.T) {
	data := ContactEmail{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "This is a test message",
		// Phone and Location are empty
	}

	body := buildEmailBody(data)

	// Required fields should be present
	if !strings.Contains(body, "Test User") {
		t.Error("Email body should contain name")
	}
	if !strings.Contains(body, "test@example.com") {
		t.Error("Email body should contain email")
	}
	if !strings.Contains(body, "This is a test message") {
		t.Error("Email body should contain message")
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"<script>", "&lt;script&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"it's", "it&#39;s"},
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("escapeHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	// Set test environment variables
	t.Setenv("SMTP_HOST", "smtp.test.com")
	t.Setenv("SMTP_PORT", "465")
	t.Setenv("SMTP_USER", "user@test.com")
	t.Setenv("SMTP_PASSWORD", "password123")
	t.Setenv("SMTP_FROM", "noreply@test.com")
	t.Setenv("SMTP_FROM_NAME", "Test App")

	config := getConfig()

	if config.Host != "smtp.test.com" {
		t.Errorf("Host = %q, want %q", config.Host, "smtp.test.com")
	}
	if config.Port != "465" {
		t.Errorf("Port = %q, want %q", config.Port, "465")
	}
	if config.User != "user@test.com" {
		t.Errorf("User = %q, want %q", config.User, "user@test.com")
	}
	if config.Password != "password123" {
		t.Errorf("Password = %q, want %q", config.Password, "password123")
	}
	if config.From != "noreply@test.com" {
		t.Errorf("From = %q, want %q", config.From, "noreply@test.com")
	}
	if config.FromName != "Test App" {
		t.Errorf("FromName = %q, want %q", config.FromName, "Test App")
	}
}
