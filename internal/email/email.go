package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type ContactEmail struct {
	Name     string
	Phone    string
	Email    string
	Location string
	Message  string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	FromName string
}

// loginAuth implements smtp.Auth for LOGIN authentication
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown LOGIN challenge: %s", fromServer)
		}
	}
	return nil, nil
}

func getConfig() SMTPConfig {
	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}
}

func SendContactEmail(data ContactEmail) error {
	config := getConfig()
	to := os.Getenv("CONTACT_EMAIL")

	if to == "" {
		return fmt.Errorf("CONTACT_EMAIL not configured")
	}

	// Log SMTP config (password obfuscated)
	fmt.Printf("SMTP Config: host=%s, port=%s, user=%s, from=%s, to=%s\n",
		config.Host, config.Port, config.User, config.From, to)

	subject := fmt.Sprintf("Nuevo contacto desde web: %s", data.Name)

	// Build HTML body
	body := buildEmailBody(data)

	// Build message
	fromHeader := config.From
	if config.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", config.FromName, config.From)
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Reply-To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		fromHeader, to, data.Email, subject, body,
	))

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// Use TLS for port 465
	if config.Port == "465" {
		return sendWithTLS(config, to, msg, addr)
	}

	// Use STARTTLS for other ports
	return sendWithSTARTTLS(config, to, msg, addr)
}

func sendWithTLS(config SMTPConfig, to string, msg []byte, addr string) error {
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}

	// Try PLAIN auth first
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}

	// Try PLAIN auth first (most common)
	plainAuth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	authErr := client.Auth(plainAuth)
	
	if authErr != nil {
		// PLAIN failed, close and try LOGIN
		client.Close()
		conn.Close()
		
		// Reconnect for LOGIN auth attempt
		conn, err = tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("TLS dial failed on retry: %w", err)
		}
		
		client, err = smtp.NewClient(conn, config.Host)
		if err != nil {
			conn.Close()
			return fmt.Errorf("SMTP client creation failed on retry: %w", err)
		}
		
		loginAuth := LoginAuth(config.User, config.Password)
		if err := client.Auth(loginAuth); err != nil {
			client.Close()
			conn.Close()
			return fmt.Errorf("SMTP auth failed (tried PLAIN: %v, LOGIN: %w)", authErr, err)
		}
	}
	
	defer client.Close()
	defer conn.Close()

	// Send email
	if err := client.Mail(config.From); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write message failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close writer failed: %w", err)
	}

	return client.Quit()
}

func sendWithSTARTTLS(config SMTPConfig, to string, msg []byte, addr string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %w", err)
	}
	defer client.Close()

	// STARTTLS
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	// Try PLAIN auth first
	plainAuth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	if err := client.Auth(plainAuth); err != nil {
		// Try LOGIN as fallback
		loginAuth := LoginAuth(config.User, config.Password)
		if err := client.Auth(loginAuth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	// Send email
	if err := client.Mail(config.From); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write message failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close writer failed: %w", err)
	}

	return client.Quit()
}

func buildEmailBody(data ContactEmail) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #000; color: #fff; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .field { margin-bottom: 15px; }
        .label { font-weight: bold; color: #3FAD4D; }
        .value { margin-top: 5px; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Nuevo Contacto</h1>
        </div>
        <div class="content">
`)

	sb.WriteString(fmt.Sprintf(`            <div class="field">
                <div class="label">Nombre:</div>
                <div class="value">%s</div>
            </div>
`, escapeHTML(data.Name)))

	if data.Phone != "" {
		sb.WriteString(fmt.Sprintf(`            <div class="field">
                <div class="label">Teléfono:</div>
                <div class="value">%s</div>
            </div>
`, escapeHTML(data.Phone)))
	}

	sb.WriteString(fmt.Sprintf(`            <div class="field">
                <div class="label">Email:</div>
                <div class="value"><a href="mailto:%s">%s</a></div>
            </div>
`, escapeHTML(data.Email), escapeHTML(data.Email)))

	if data.Location != "" {
		sb.WriteString(fmt.Sprintf(`            <div class="field">
                <div class="label">Ubicación:</div>
                <div class="value">%s</div>
            </div>
`, escapeHTML(data.Location)))
	}

	sb.WriteString(fmt.Sprintf(`            <div class="field">
                <div class="label">Mensaje:</div>
                <div class="value">%s</div>
            </div>
`, escapeHTML(data.Message)))

	sb.WriteString(`        </div>
        <div class="footer">
            Este mensaje fue enviado desde el formulario de contacto de la web.
        </div>
    </div>
</body>
</html>`)

	return sb.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
