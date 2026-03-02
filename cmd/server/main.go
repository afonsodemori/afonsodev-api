package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/afonsodemori/afonsodev-api/internal/challenge"
	"github.com/afonsodemori/afonsodev-api/internal/config"
	"github.com/afonsodemori/afonsodev-api/internal/contact"
	"github.com/afonsodemori/afonsodev-api/internal/email"
	apphttp "github.com/afonsodemori/afonsodev-api/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Wire dependencies
	recaptcha := challenge.NewRecaptchaVerifier(cfg.RecaptchaSecret)
	turnstile := challenge.NewTurnstileVerifier(cfg.TurnstileSecret)
	emailClient := email.NewResendClient(cfg.ResendAPIKey)

	contactService := contact.NewService(
		contact.ServiceConfig{
			IsDevelopment: cfg.IsDevelopment(),
			ContactFrom:   cfg.ContactFrom,
			ContactTo:     cfg.ContactTo,
		},
		recaptcha,
		turnstile,
		emailClient,
	)

	contactHandler := contact.NewHandler(contactService)

	router := apphttp.NewRouter(cfg.AllowedOrigins, contactHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
