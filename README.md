# afonsodev-api

Backend API for [afonso.dev](https://afonso.dev), written in Go.

## What it does

Handles the contact form on my personal website — validates submissions, verifies bot-protection challenges, and sends emails via [Resend](https://resend.com).

## Stack

- **Go** — no frameworks, stdlib `net/http` only
- **Cloudflare Turnstile** and **Google reCAPTCHA** for bot protection
- **Resend** for transactional email
- **Docker** — multi-arch image (`linux/amd64`, `linux/arm64`) built from scratch
- **Cloudflare Tunnel** for secure dev exposure without opening ports

## Highlights

- Clean vertical slice architecture: each feature owns its model, service, handler, and errors
- Interfaces over concrete types — `ChallengeVerifier` and `EmailSender` keep the business logic decoupled and testable
- No third-party HTTP router or DI container — intentionally minimal
- Dev container setup with VS Code for a one-command dev environment

## Running locally

```bash
cp .env.example .env  # fill in your secrets
make run
```

The server starts on port `8080` by default. Set `ENV=development` to skip actual email sending.

## License

[MIT](LICENSE)
