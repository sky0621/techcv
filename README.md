# techcv
Curriculum Vitae – Technical Proficiencies of a Web Engineer

## Backend Tooling

From `services/manager/backend`, run:

- `make generate` – regenerate Echo-compatible handlers and types from `docs/openapi.yaml`.
- `VERIFICATION_URL_BASE` – optional, defaults to `http://localhost:5173/auth/verify`; used by the manager API when composing verification links in registration emails.
