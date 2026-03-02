package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port            string
	Env             string
	AllowedOrigins  []string
	RecaptchaSecret string
	TurnstileSecret string
	ResendAPIKey    string
	ContactFrom     string
	ContactTo       string
}

func (c Config) IsDevelopment() bool {
	return c.Env == "development"
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:            port,
		Env:             os.Getenv("ENV"),
		RecaptchaSecret: os.Getenv("RECAPTCHA_SECRET"),
		TurnstileSecret: os.Getenv("TURNSTILE_SECRET"),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ContactFrom:     os.Getenv("CONTACT_FROM"),
		ContactTo:       os.Getenv("CONTACT_TO"),
	}

	if origins := os.Getenv("ALLOWED_ORIGIN"); origins != "" {
		cfg.AllowedOrigins = strings.Split(origins, ",")
	}

	if !cfg.IsDevelopment() {
		if cfg.ResendAPIKey == "" || cfg.ContactFrom == "" || cfg.ContactTo == "" {
			return nil, fmt.Errorf("RESEND_API_KEY, CONTACT_FROM, and CONTACT_TO must be set in non-development environments")
		}
	}

	return cfg, nil
}
