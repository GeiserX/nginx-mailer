# Getting started

nginx-mailer is one container: nginx serves your static site on port 80 and proxies `/api/` to a small Go service that sends the contact form by SMTP. The image is built for `linux/amd64`.

## Docker

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
  drumsergio/nginx-mailer:1.0.0
```

## Docker Compose

The same setup as a compose file (also in [`docker-compose.example.yml`](https://github.com/GeiserX/nginx-mailer/blob/main/docker-compose.example.yml)):

```yaml
services:
  website:
    image: drumsergio/nginx-mailer:1.0.0
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

## Check that it works

`curl http://localhost/health` returns `OK`, and your site is on http://localhost. Every variable is in [Configuration](configuration.md); the form API is in [Usage](usage.md).
