package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/GeiserX/nginx-mailer/internal/email"
	"github.com/GeiserX/nginx-mailer/internal/turnstile"
)

type ContactRequest struct {
	Name      string `json:"nombre"`
	Phone     string `json:"telefono"`
	Email     string `json:"email"`
	Location  string `json:"ubicacion"`
	Message   string `json:"mensaje"`
	Turnstile string `json:"cf-turnstile-response"`
}

type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ContactRequest

	// Parse form data or JSON
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Error parsing request", http.StatusBadRequest)
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Multipart form data (from FormData in JS)
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
			sendError(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}
		req = ContactRequest{
			Name:      r.FormValue("nombre"),
			Phone:     r.FormValue("telefono"),
			Email:     r.FormValue("email"),
			Location:  r.FormValue("ubicacion"),
			Message:   r.FormValue("mensaje"),
			Turnstile: r.FormValue("cf-turnstile-response"),
		}
	} else {
		// URL-encoded form data
		if err := r.ParseForm(); err != nil {
			sendError(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		req = ContactRequest{
			Name:      r.FormValue("nombre"),
			Phone:     r.FormValue("telefono"),
			Email:     r.FormValue("email"),
			Location:  r.FormValue("ubicacion"),
			Message:   r.FormValue("mensaje"),
			Turnstile: r.FormValue("cf-turnstile-response"),
		}
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Message == "" {
		sendError(w, "Nombre, email y mensaje son obligatorios", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !strings.Contains(req.Email, "@") {
		sendError(w, "Email inválido", http.StatusBadRequest)
		return
	}

	// Verify Turnstile CAPTCHA
	if err := turnstile.Verify(req.Turnstile, r.RemoteAddr); err != nil {
		log.Printf("Turnstile verification failed: %v", err)
		sendError(w, "Verificación de seguridad fallida", http.StatusForbidden)
		return
	}

	// Send email
	emailData := email.ContactEmail{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Location: req.Location,
		Message:  req.Message,
	}

	if err := email.SendContactEmail(emailData); err != nil {
		log.Printf("Failed to send email: %v", err)
		sendError(w, "Error al enviar el mensaje", http.StatusInternalServerError)
		return
	}

	log.Printf("Contact form submitted successfully from %s <%s>", req.Name, req.Email)

	if err := json.NewEncoder(w).Encode(ContactResponse{
		Success: true,
		Message: "Mensaje enviado correctamente. Nos pondremos en contacto contigo pronto.",
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ContactResponse{
		Success: false,
		Message: message,
	}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
