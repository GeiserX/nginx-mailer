package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestContactHandler_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		body     map[string]string
		wantCode int
	}{
		{
			name:     "missing all fields",
			body:     map[string]string{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing email",
			body:     map[string]string{"nombre": "Test", "mensaje": "Hello"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing name",
			body:     map[string]string{"email": "test@test.com", "mensaje": "Hello"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing message",
			body:     map[string]string{"nombre": "Test", "email": "test@test.com"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid email",
			body:     map[string]string{"nombre": "Test", "email": "invalid", "mensaje": "Hello"},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			ContactHandler(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("ContactHandler() status = %v, want %v", rr.Code, tt.wantCode)
			}

			var resp ContactResponse
			json.NewDecoder(rr.Body).Decode(&resp)
			if resp.Success {
				t.Errorf("ContactHandler() success = true, want false")
			}
		})
	}
}

func TestContactHandler_FormData(t *testing.T) {
	form := url.Values{}
	form.Set("nombre", "Test User")
	form.Set("email", "test@test.com")
	form.Set("mensaje", "Hello World")

	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	ContactHandler(rr, req)

	// Without SMTP config, it will fail at email sending, but validation passes
	// This tests that form parsing works correctly
	if rr.Code == http.StatusBadRequest {
		t.Errorf("ContactHandler() should not return BadRequest for valid form data")
	}
}

func TestContactHandler_JSON(t *testing.T) {
	body := map[string]string{
		"nombre":  "Test User",
		"email":   "test@test.com",
		"mensaje": "Hello World",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/contact", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ContactHandler(rr, req)

	// Without SMTP config, it will fail at email sending, but validation passes
	if rr.Code == http.StatusBadRequest {
		t.Errorf("ContactHandler() should not return BadRequest for valid JSON data")
	}
}
