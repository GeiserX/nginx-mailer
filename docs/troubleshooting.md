# Troubleshooting

Both nginx and the mail service log to the container output: `docker logs <container>`.

## The form answers "Error al enviar el mensaje" (HTTP 500)

The SMTP send failed; the log line `Failed to send email:` says why. Check that `SMTP_PORT` matches the server: `465` is TLS from the first byte, `587` is STARTTLS, and a mismatch fails at the TLS dial. Check `SMTP_USER` and `SMTP_PASSWORD`, and that `CONTACT_EMAIL` is set (`CONTACT_EMAIL not configured` in the log).

## The form answers "Verificación de seguridad fallida" (HTTP 403)

`CLOUDFLARE_TURNSTILE_SECRET_KEY` is set and the Turnstile check failed: the form sent no `cf-turnstile-response`, or the site key on the page does not belong to the secret key. The log line `Turnstile verification failed:` gives the reason. An empty variable turns the check off, which is only for local testing: without it, anyone who finds the endpoint can make the container send mail.

## The form answers "Nombre, email y mensaje son obligatorios" (HTTP 400)

The request lacks `nombre`, `email` or `mensaje`, or the field names differ from the ones in [Usage](usage.md).

## Reporting a bug

Open an issue at https://github.com/GeiserX/nginx-mailer/issues with the image tag, the request you sent, the response, and the container log. Remove the SMTP password and any names, email addresses or messages from past submissions first: the log records who submitted the form.
