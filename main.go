package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type SendEmailRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Details any    `json:"details,omitempty"`
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://afonso.dev", http.StatusTemporaryRedirect)
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGIN")
	if allowedOriginsStr != "" {
		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		origin := r.Header.Get("Origin")
		for _, ao := range allowedOrigins {
			if origin == ao {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
	} else {
		// TODO: Remove this for production
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Message: "Method Not Allowed"})
		return
	}

	var req SendEmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Invalid request payload"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Message: "contact.form.missing_fields"})
		return
	}

	if req.Token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "contact.form.captcha.missing"})
		return
	}

	recaptchaSuccess, recaptchaErrors, err := VerifyRecaptcha(req.Token)
	if err != nil {
		log.Printf("reCAPTCHA verification error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "reCAPTCHA verification failed"})
		return
	}

	if !recaptchaSuccess {
		log.Printf("reCAPTCHA verification unsuccessful: %v", recaptchaErrors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendEmailResponse{
			Success: false,
			Error:   "contact.form.captcha.invalid",
			Details: recaptchaErrors,
		})
		return
	}

	resendApiResponse, err := SendEmail(req.Name, req.Email, req.Subject, req.Message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{
			Success: false,
			Error:   "Failed to send email",
			Details: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SendEmailResponse{
		Success: true,
		Message: "contact.form.success",
		Details: resendApiResponse,
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/send-email", sendEmailHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
