# afonsodev-api

A Go-based microservice that handles contact form submissions for the [afonso.dev](https://afonso.dev) website. It validates submissions via Google reCAPTCHA or Cloudflare Turnstile and sends emails using the Resend API.

## Project Overview

- **Core Functionality:** Validates and processes contact form submissions.
- **Main Technologies:**
  - **Language:** Go (1.25.6)
  - **Verification:** Google reCAPTCHA v3, Cloudflare Turnstile
  - **Email Provider:** [Resend](https://resend.com)
  - **Infrastructure:** Cloudflare Tunnel for development, Multi-platform Docker for deployment.
- **Architecture:** Simple HTTP server using the Go standard library.

## Key Files & Directories

- `main.go`: Entry point, handles routing and HTTP requests.
- `email.go`: Integration logic for the Resend API.
- `captcha.go`: Logic for verifying Google reCAPTCHA tokens.
- `turnstile.go`: Logic for verifying Cloudflare Turnstile tokens.
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

- **CORS:** Controlled via the `ALLOWED_ORIGIN` environment variable.
- **Verification:** Default challenger is `captcha` (reCAPTCHA) if not specified.
- **Logging:** Errors are logged to standard output with relevant context.
- **Deployment:** Uses Docker `buildx` for multi-platform support (`amd64` and `arm64`).
