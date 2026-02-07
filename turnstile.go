package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

func VerifyTurnstile(token string) (bool, []string, error) {
	turnstileVerifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	turnstileSecret := os.Getenv("TURNSTILE_SECRET")

	if turnstileSecret == "" {
		return false, nil, fmt.Errorf("TURNSTILE_SECRET environment variable not set")
	}

	requestBody := map[string]string{
		"secret":   turnstileSecret,
		"response": token,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, nil, fmt.Errorf("error marshalling turnstile request body: %w", err)
	}

	req, err := http.NewRequest("POST", turnstileVerifyURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, nil, fmt.Errorf("error creating turnstile request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("error making turnstile request to %s: %w", turnstileVerifyURL, err)
	}
	defer res.Body.Close()

	var turnstileResponseData TurnstileResponse
	if err := json.NewDecoder(res.Body).Decode(&turnstileResponseData); err != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Error decoding turnstile response: %v, raw body: %s", err, string(bodyBytes))
		return false, nil, fmt.Errorf("error decoding turnstile response: %w", err)
	}

	return turnstileResponseData.Success, turnstileResponseData.ErrorCodes, nil
}
