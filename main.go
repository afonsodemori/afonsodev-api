package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

type RecaptchaResponse struct {
	Success bool     `json:"success"`
	Score   float64  `json:"score"`
	Action  string   `json:"action"`
	Errors  []string `json:"error-codes"`
}

type ResendEmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	ReplyTo string `json:"reply_to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Added Authorization header

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

	recaptchaVerifyUrl := "https://www.google.com/recaptcha/api/siteverify"
	recaptchaSecret := os.Getenv("RECAPTCHA_SECRET")

	if recaptchaSecret == "" {
		log.Println("RECAPTCHA_SECRET environment variable not set.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Server configuration error: RECAPTCHA_SECRET not set"})
		return
	}

	recaptchaForm := url.Values{}
	recaptchaForm.Add("secret", recaptchaSecret)
	recaptchaForm.Add("response", req.Token)

	recaptchaRes, err := http.PostForm(recaptchaVerifyUrl, recaptchaForm)
	if err != nil {
		log.Printf("Error making reCAPTCHA request to %s: %v", recaptchaVerifyUrl, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "reCAPTCHA verification failed"})
		return
	}
	defer recaptchaRes.Body.Close()

	var recaptchaResponseData RecaptchaResponse
	if err := json.NewDecoder(recaptchaRes.Body).Decode(&recaptchaResponseData); err != nil {
		log.Printf("Error decoding reCAPTCHA response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "reCAPTCHA verification failed (decode)"})
		return
	}

	if !recaptchaResponseData.Success {
		log.Printf("reCAPTCHA verification unsuccessful: %v", recaptchaResponseData.Errors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendEmailResponse{
			Success: false,
			Error:   "contact.form.captcha.invalid",
			Details: recaptchaResponseData.Errors,
		})
		return
	}

	resendApiUrl := "https://api.resend.com/emails"
	resendApiKey := os.Getenv("RESEND_API_KEY")
	contactFrom := os.Getenv("CONTACT_FROM")
	contactTo := os.Getenv("CONTACT_TO")

	if resendApiKey == "" || contactFrom == "" || contactTo == "" {
		log.Println("Resend API environment variables (RESEND_API_KEY, CONTACT_FROM, CONTACT_TO) not set.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Server configuration error: Resend API keys not set"})
		return
	}

	resendReqBody := ResendEmailRequest{
		From:    fmt.Sprintf("%s - via afonso.dev <%s>", req.Name, contactFrom),
		To:      contactTo,
		ReplyTo: req.Email,
		Subject: req.Subject,
		Text:    req.Message,
	}

	resendReqBodyBytes, err := json.Marshal(resendReqBody)
	if err != nil {
		log.Printf("Error marshalling Resend email request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Failed to prepare email"})
		return
	}

	resendClient := &http.Client{}
	resendApiReq, err := http.NewRequest("POST", resendApiUrl, bytes.NewBuffer(resendReqBodyBytes))
	if err != nil {
		log.Printf("Error creating Resend API request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Failed to send email (request creation)"})
		return
	}
	resendApiReq.Header.Set("Content-Type", "application/json")
	resendApiReq.Header.Set("Authorization", "Bearer "+resendApiKey)

	resendApiRes, err := resendClient.Do(resendApiReq)
	if err != nil {
		log.Printf("Error making Resend API request to %s: %v", resendApiUrl, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendEmailResponse{Success: false, Error: "Failed to send email (API call)"})
		return
	}
	defer resendApiRes.Body.Close()

	resendResponseBody, _ := io.ReadAll(resendApiRes.Body)
	if resendApiRes.StatusCode != http.StatusOK {
		log.Printf("Resend API returned non-200 status: %d, response: %s", resendApiRes.StatusCode, string(resendResponseBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(SendEmailResponse{
			Success: false,
			Error:   "Failed to send email (Resend API error)",
			Details: json.RawMessage(resendResponseBody), // Pass raw JSON for details
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SendEmailResponse{
		Success: true,
		Message: "contact.form.success",
		Details: json.RawMessage(resendResponseBody),
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
