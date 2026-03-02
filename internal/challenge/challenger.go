package challenge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Result holds the outcome of a challenge verification.
type Result struct {
	Success bool
	Errors  []string
}

// Verifier verifies a challenge token.
type Verifier interface {
	Verify(ctx context.Context, token string) (*Result, error)
}

// --- reCAPTCHA ---

type recaptchaResponse struct {
	Success bool     `json:"success"`
	Score   float64  `json:"score"`
	Action  string   `json:"action"`
	Errors  []string `json:"error-codes"`
}

type RecaptchaVerifier struct {
	secret string
}

func NewRecaptchaVerifier(secret string) *RecaptchaVerifier {
	return &RecaptchaVerifier{secret: secret}
}

func (v *RecaptchaVerifier) Verify(ctx context.Context, token string) (*Result, error) {
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	form := url.Values{}
	form.Add("secret", v.secret)
	form.Add("response", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Creating recaptcha request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Recaptcha request to %s: %w", verifyURL, err)
	}
	defer res.Body.Close()

	var data recaptchaResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Error decoding reCAPTCHA response: %v, raw body: %s", err, string(bodyBytes))
		return nil, fmt.Errorf("Decoding recaptcha response: %w", err)
	}

	return &Result{Success: data.Success, Errors: data.Errors}, nil
}

// --- Turnstile ---

type turnstileRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type turnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

type TurnstileVerifier struct {
	secret string
}

func NewTurnstileVerifier(secret string) *TurnstileVerifier {
	return &TurnstileVerifier{secret: secret}
}

func (v *TurnstileVerifier) Verify(ctx context.Context, token string) (*Result, error) {
	verifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	body := turnstileRequest{
		Secret:   v.secret,
		Response: token,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Marshalling turnstile request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Creating turnstile request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Turnstile request to %s: %w", verifyURL, err)
	}
	defer res.Body.Close()

	var data turnstileResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Error decoding turnstile response: %v, raw body: %s", err, string(bodyBytes))
		return nil, fmt.Errorf("Decoding turnstile response: %w", err)
	}

	return &Result{Success: data.Success, Errors: data.ErrorCodes}, nil
}
