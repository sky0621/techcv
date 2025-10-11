# Repository Guidelines

## Project Structure & Module Organization
The backend lives in `services`, a standalone Go module with the HTTP entrypoint under `services/cmd/api`. Layered packages are in `services/internal` (`domain`, `infrastructure`, `interface`, `usecase`)—keep new business logic in the relevant layer to preserve separation of concerns. Shared documentation belongs in `docs`, while deployment manifests and scripts belong in `infra`.

## Build, Test, and Development Commands
Run tooling from the `services` directory:
- `make run`: starts the Echo API locally via `go run ./cmd/api`, respecting the `PORT` environment variable.
- `make build`: compiles `bin/api`, the production binary built from `cmd/api`.
- `make test`: executes `go test ./...` across every package.
- `make lint`: runs `gofmt` on all Go sources.
- `make tidy`: refreshes module dependencies and prunes unused ones.

## Coding Style & Naming Conventions
Follow idiomatic Go style—use tabs, keep line length reasonable, and rely on `gofmt` for layout. Package names should be short and lower-case (`handler`, `logger`); exported structs and interfaces use PascalCase, private helpers use camelCase. Prefer constructor functions like `NewHealthHandler` for dependency injection and keep HTTP handlers in `internal/interface/http` alongside their middleware.

## Testing Guidelines
Add `_test.go` files alongside the code they exercise and use the standard `testing` package. Group table-driven cases when inputs vary, and stub external dependencies via interfaces defined in `internal`. Aim for meaningful coverage on new code paths and ensure `make test` passes cleanly before pushing.

## Commit & Pull Request Guidelines
Write commits in the imperative mood with concise subjects (`Add health handler timeout logging`). The current history favors descriptive messages that enumerate impacted paths; mirror that clarity and note structural moves when applicable. Pull requests must explain the user-facing change, list testing evidence (`make test` output is sufficient), and link tracking issues. Add screenshots or curl transcripts when you alter HTTP responses.

## Environment & Configuration
The API reads configuration from environment variables; currently only `PORT` is required, defaulting to `8080`. Centralized logging uses Zap (`internal/infrastructure/logger`), so prefer structured fields (`zap.String`, etc.) over string concatenation. Update the README whenever configuration inputs change so deployment manifests in `infra` stay aligned.
