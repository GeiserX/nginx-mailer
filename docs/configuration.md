# Configuration

All settings are environment variables on the container.

| Variable | Required | Description |
|----------|----------|-------------|
| `SMTP_HOST` | Yes | SMTP server hostname |
| `SMTP_PORT` | Yes | SMTP port: `465` connects with TLS from the start; any other port (usually `587`) uses STARTTLS |
| `SMTP_USER` | Yes | SMTP username |
| `SMTP_PASSWORD` | Yes | SMTP password |
| `SMTP_FROM` | Yes | From email address |
| `SMTP_FROM_NAME` | No | From display name |
| `CONTACT_EMAIL` | Yes | Recipient email for contact forms |
| `CLOUDFLARE_TURNSTILE_SECRET_KEY` | No | Turnstile secret key; when empty, the CAPTCHA check is skipped |

Mount your site at `/usr/share/nginx/html` (read-only is fine). nginx serves it on port 80.
