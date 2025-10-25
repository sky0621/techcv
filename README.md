# techcv
Curriculum Vitae – Technical Proficiencies of a Web Engineer

## Backend Tooling

From `services/manager/backend`, run:

- `make generate` – regenerate Echo-compatible handlers and types from `docs/openapi.yaml`.

### Environment Variables

Configure the manager backend by copying `.env.example` to `.env` (or exporting the variables manually). Required values include:

- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URI` – Google OAuth 2.0 credentials.
- `SESSION_SECRET` – random string used to sign session tokens.
- `COOKIE_DOMAIN`, `COOKIE_SECURE` – HTTP cookie scope and security mode.
- `REDIS_ADDR` (with optional `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB`) – session store connection.
- `VERIFICATION_URL_BASE` – optional, defaults to `http://localhost:5173/auth/verify`; used by the manager API when composing verification links in registration emails.
