# Usage

## API

### POST /api/contact

Accepts JSON, `multipart/form-data` or a URL-encoded form with these fields. `nombre`, `email` and `mensaje` are required; `cf-turnstile-response` is required when `CLOUDFLARE_TURNSTILE_SECRET_KEY` is set.

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

The answer is JSON: `{"success": true, "message": "..."}` on success, `success: false` with a message otherwise. The validation and result messages are in Spanish; request parsing errors are in English.

### GET /health

Returns `200 OK` for health checks.

## HTML form example

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
