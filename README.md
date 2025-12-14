# nginx-mailer

Lightweight Docker image based on nginx:alpine that serves static websites and handles contact form submissions via SMTP.

## Features

- 🚀 Based on `nginx:1.29-alpine` (~35MB total)
- 📧 Built-in contact form API with SMTP support
- 🔒 Cloudflare Turnstile CAPTCHA integration
- 🐳 Single container (nginx + Go API)
- 🏗️ Multi-arch: `linux/amd64` and `linux/arm64`

## Quick Start

```bash
docker run -d \
  -p 80:80 \
  -v /path/to/site:/usr/share/nginx/html:ro \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=465 \
  -e SMTP_USER=noreply@example.com \
  -e SMTP_PASSWORD=your-password \
  -e SMTP_FROM=noreply@example.com \
  -e CONTACT_EMAIL=you@example.com \
  drumsergio/nginx-mailer:latest
```

## Docker Compose

```yaml
services:
  website:
    image: drumsergio/nginx-mailer:latest
    ports:
      - "80:80"
    volumes:
      - ./website:/usr/share/nginx/html:ro
    environment:
      - SMTP_HOST=smtp.example.com
      - SMTP_PORT=465
      - SMTP_USER=noreply@example.com
      - SMTP_PASSWORD=your-password
      - SMTP_FROM=noreply@example.com
      - SMTP_FROM_NAME=My Website
      - CONTACT_EMAIL=you@example.com
      - CLOUDFLARE_TURNSTILE_SECRET_KEY=  # Optional
    restart: unless-stopped
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SMTP_HOST` | Yes | SMTP server hostname |
| `SMTP_PORT` | Yes | SMTP port (465 for SSL, 587 for STARTTLS) |
| `SMTP_USER` | Yes | SMTP username |
| `SMTP_PASSWORD` | Yes | SMTP password |
| `SMTP_FROM` | Yes | From email address |
| `SMTP_FROM_NAME` | No | From display name |
| `CONTACT_EMAIL` | Yes | Recipient email for contact forms |
| `CLOUDFLARE_TURNSTILE_SECRET_KEY` | No | Turnstile secret key |

## API

### POST /api/contact

```json
{
  "nombre": "John Doe",
  "email": "john@example.com",
  "telefono": "+1234567890",
  "ubicacion": "City",
  "mensaje": "Hello...",
  "cf-turnstile-response": "token"
}
```

### GET /health

Returns `200 OK` for health checks.

## HTML Form Example

```html
<form action="/api/contact" method="POST">
  <input type="text" name="nombre" placeholder="Name" required>
  <input type="email" name="email" placeholder="Email" required>
  <input type="tel" name="telefono" placeholder="Phone">
  <textarea name="mensaje" placeholder="Message" required></textarea>
  
  <!-- Optional: Cloudflare Turnstile -->
  <div class="cf-turnstile" data-sitekey="your-site-key"></div>
  <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
  
  <button type="submit">Send</button>
</form>
```

## License

MIT
