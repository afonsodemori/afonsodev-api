# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make run              # Run the server locally (loads .env automatically)
go test ./...         # Run all tests
go test ./internal/contact/...  # Run tests for a specific package
make docker-build     # Build multi-arch Docker image
make docker-push      # Build and push Docker image
make tunnel-run       # Run Cloudflare tunnel (dev exposure)
```

The `.env` file is auto-loaded by the Makefile. Copy `.env.example` to `.env` to get started.

## Architecture

This is a Go HTTP API (`github.com/afonsodemori/afonsodev-api`) with no external framework dependencies — only stdlib `net/http`.

**Entry point:** `cmd/server/main.go` wires all dependencies manually (no DI framework) and starts the server.

**Package structure:**

- `internal/config` — loads all config from env vars; exposes `Commit`/`BuiltAt` build-time variables injected via `-ldflags` in the Docker build
- `internal/http` — router, CORS middleware, and health/redirect handlers
- `internal/challenge` — `RecaptchaVerifier` and `TurnstileVerifier`, both implementing the `Verifier` interface; selected per-request via the `challenger` field
- `internal/contact` — full vertical slice: model, errors, service (business logic), handler (HTTP); `Service` depends on `ChallengeVerifier` and `EmailSender` interfaces
- `internal/email` — `ResendClient` implementing `EmailSender`

**Key behavior:**

- In `development` env (`ENV=development`), challenge verification still runs but email sending is skipped
- `RESEND_API_KEY`, `CONTACT_FROM`, and `CONTACT_TO` are required in non-development environments; startup fails fast if missing
- `/send-email` is a deprecated alias for `/contact` (both handled by `contact.Handler.HandleSendEmail`)
- Error messages returned to clients use i18n-style dot-notation keys (e.g., `contact.form.missing_fields`) that the frontend translates
