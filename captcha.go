package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type RecaptchaResponse struct {
	Success bool     `json:"success"`
	Score   float64  `json:"score"`
	Action  string   `json:"action"`
	Errors  []string `json:"error-codes"`
}

func VerifyRecaptcha(token string) (bool, []string, error) {
	recaptchaVerifyUrl := "https://www.google.com/recaptcha/api/siteverify"
	recaptchaSecret := os.Getenv("RECAPTCHA_SECRET")

	if recaptchaSecret == "" {
		return false, nil, fmt.Errorf("RECAPTCHA_SECRET environment variable not set")
	}

	recaptchaForm := url.Values{}
	recaptchaForm.Add("secret", recaptchaSecret)
	recaptchaForm.Add("response", token)

	res, err := http.PostForm(recaptchaVerifyUrl, recaptchaForm)
	if err != nil {
		return false, nil, fmt.Errorf("error making reCAPTCHA request to %s: %w", recaptchaVerifyUrl, err)
	}
	defer res.Body.Close()

	var recaptchaResponseData RecaptchaResponse
	if err := json.NewDecoder(res.Body).Decode(&recaptchaResponseData); err != nil {
		bodyBytes, _ := io.ReadAll(res.Body) // Read body for logging if decode fails
		log.Printf("Error decoding reCAPTCHA response: %v, raw body: %s", err, string(bodyBytes))
		return false, nil, fmt.Errorf("error decoding reCAPTCHA response: %w", err)
	}

	return recaptchaResponseData.Success, recaptchaResponseData.Errors, nil
}
