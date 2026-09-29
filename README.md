<p align="center"><img src="https://raw.githubusercontent.com/GeiserX/nginx-mailer/main/docs/images/banner.svg" alt="nginx-mailer banner" width="900"/></p>

<h1 align="center">nginx-mailer</h1>

<p align="center">
  <a href="https://github.com/GeiserX/nginx-mailer/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/GeiserX/nginx-mailer/ci.yml?label=CI" alt="CI"/></a>
  <a href="https://github.com/GeiserX/nginx-mailer/blob/main/LICENSE"><img src="https://img.shields.io/github/license/GeiserX/nginx-mailer" alt="License"/></a>
  <a href="https://hub.docker.com/r/drumsergio/nginx-mailer"><img src="https://img.shields.io/docker/pulls/drumsergio/nginx-mailer" alt="Docker Pulls"/></a>
</p>

<p align="center"><strong>Docker image based on nginx:alpine that serves a static site and sends its contact form by SMTP, with optional Cloudflare Turnstile.</strong></p>

---

## Features

- Serves a static site from a mounted folder with nginx on port 80.
- Contact form endpoint `POST /api/contact` that emails each submission to one address.
- SMTP over TLS (port 465) or STARTTLS (other ports).
- Optional Cloudflare Turnstile CAPTCHA check.
- One container: nginx and the Go mail service run under supervisord.
- `linux/amd64` image on Docker Hub.

## Quick start

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

The site is served on port 80, with `/api/contact` and `/health` behind nginx. Compose and every variable are in [Getting started](https://github.com/GeiserX/nginx-mailer/blob/main/docs/getting-started.md).

## Documentation

- [Getting started](https://github.com/GeiserX/nginx-mailer/blob/main/docs/getting-started.md): `docker run`, Docker Compose, and the health check
- [Configuration](https://github.com/GeiserX/nginx-mailer/blob/main/docs/configuration.md): every environment variable
- [Usage](https://github.com/GeiserX/nginx-mailer/blob/main/docs/usage.md): the contact API and an HTML form example
- [Troubleshooting](https://github.com/GeiserX/nginx-mailer/blob/main/docs/troubleshooting.md): SMTP and Turnstile failures, and what to put in a bug report

## License

[GPL-3.0-or-later](https://github.com/GeiserX/nginx-mailer/blob/main/LICENSE)
