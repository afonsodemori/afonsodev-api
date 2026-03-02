# afonsodev-api

A Go-based microservice that handles contact form submissions for the [afonso.dev](https://afonso.dev) website. It validates submissions via Google reCAPTCHA or Cloudflare Turnstile and sends emails using the Resend API.

## Project Overview

- **Core Functionality:** Validates and processes contact form submissions.
- **Main Technologies:**
  - **Language:** Go (1.25.6)
  - **Verification:** Google reCAPTCHA v3, Cloudflare Turnstile
  - **Email Provider:** [Resend](https://resend.com)
  - **Infrastructure:** Cloudflare Tunnel for development, Multi-platform Docker for deployment.
- **Architecture:** Layered HTTP server using the Go standard library, organized by feature domain.

## Key Files & Directories

- `cmd/server/main.go`: Application bootstrap (config, dependency wiring, server start).
- `internal/config/config.go`: Centralized configuration parsing from environment variables.
- `internal/http/router.go`: Global HTTP router, CORS middleware, and the redirect handler.
- `internal/contact/`: Core feature package for the contact form.
  - `model.go`: Domain entities (request/response types).
  - `service.go`: Business logic (validation, challenge verification dispatch, email sending).
  - `handler.go`: HTTP handlers; maps domain errors to HTTP status codes.
  - `errors.go`: Domain-level sentinel errors.
- `internal/challenge/challenger.go`: `Verifier` interface with reCAPTCHA and Turnstile implementations.
- `internal/email/resend.go`: `Sender` interface with Resend API implementation.
- `Makefile`: Automation for running, tunneling, and building Docker images.
- `docker/`: Contains the production Dockerfile.
- `.devcontainer/`: Configuration for the development container and Cloudflare Tunnel.

## API Endpoints

### `GET /`

Redirects temporarily to `https://afonso.dev`.

### `POST /send-email`

Validates the request and sends an email via Resend.

- **Request Body (JSON):**
  ```json
  {
    "name": "string",
    "email": "string",
    "subject": "string",
    "message": "string",
    "token": "string",
    "challenger": "captcha" | "turnstile"
  }
  ```
- **Responses:**
  - `200 OK`: Success, email sent.
  - `400 Bad Request`: Missing fields, invalid token, or unknown challenger.
  - `500 Internal Server Error`: Server-side failures (API errors).

## Development Guide

### Prerequisites

- Go 1.25+ installed locally or via Dev Container.
- [cloudflared](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/install-and-setup/installation/) for tunneling.
- Docker for containerized builds.

### Environment Variables

Create a `.env` file based on `.env.example`:

- `ENV`: Set to `development` to disable actual email sending.
- `ALLOWED_ORIGIN`: Comma-separated list of origins for CORS.
- `RECAPTCHA_SECRET`: Google reCAPTCHA secret key.
- `TURNSTILE_SECRET`: Cloudflare Turnstile secret key.
- `RESEND_API_KEY`: API key for Resend.
- `CONTACT_FROM`: Verified sender email in Resend.
- `CONTACT_TO`: Destination email for form submissions.

### Common Commands

- **Run Locally:** `make run`
- **Initialize Cloudflare Tunnel:** `make tunnel-init`
- **Run Cloudflare Tunnel:** `make tunnel-run`
- **Build Multi-Platform Docker Image:** `make docker-build`

## Development Conventions

### Project Structure

- Follows the standard Go project layout.
- Uses `internal/` for private application code to prevent external importing.
- Organized by **feature (domain)** rather than technical layers. Each feature package (e.g., `internal/contact`) contains its own models, service logic, and HTTP handlers.
- `internal/http` handles routing and global middleware (CORS).
- Dependencies flow inward: **[HTTP Handler] → [Service] → [Repository/External Client]**.

### Coding Style

- Adhere to idiomatic Go (Effective Go).
- **Interfaces:** Define interfaces in the consuming package. Accept interfaces as parameters. Return concrete types from constructors.
- **Context:** Pass `context.Context` explicitly through all layers.
- **Config:** All configuration is parsed once at startup into a `config.Config` struct — no `os.Getenv` calls outside `internal/config`.

### Guidelines

- **Idiomatic Go:** Use the standard library where possible.
- **Clear Boundaries:**
  - Handlers know HTTP.
  - Services know business rules.
  - External clients (email, challenge) know their respective APIs.
  - Models remain transport-agnostic.
- **Error Strategy:**
  - Define domain-level sentinel errors in the feature package (e.g., `ErrMissingFields`, `ErrInvalidToken`).
  - Services return only domain errors — never HTTP status codes.
  - Handlers translate domain errors to HTTP responses (e.g., `ErrMissingFields` → 400).
- **CORS:** Controlled via the `ALLOWED_ORIGIN` environment variable, applied as global middleware.
- **Verification:** Default challenger is `captcha` (reCAPTCHA) if not specified.
- **Logging:** Errors are logged to standard output with relevant context.
- **Deployment:** Uses Docker `buildx` for multi-platform support (`amd64` and `arm64`).
- **Keep It Simple:**
  - No premature abstractions.
  - Use interfaces for external services to facilitate testing and decouple layers.
