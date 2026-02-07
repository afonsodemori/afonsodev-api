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

type ResendEmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	ReplyTo string `json:"reply_to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func SendEmail(name, email, subject, message string) (any, error) {
	if os.Getenv("ENV") == "development" {
		log.Printf("Development mode. Email was not sent")
		return nil, nil
	}

	resendApiUrl := "https://api.resend.com/emails"
	resendApiKey := os.Getenv("RESEND_API_KEY")
	contactFrom := os.Getenv("CONTACT_FROM")
	contactTo := os.Getenv("CONTACT_TO")

	if resendApiKey == "" || contactFrom == "" || contactTo == "" {
		return nil, fmt.Errorf("Resend API environment variables (RESEND_API_KEY, CONTACT_FROM, CONTACT_TO) not set")
	}

	resendReqBody := ResendEmailRequest{
		From:    fmt.Sprintf("%s - via afonso.dev <%s>", name, contactFrom),
		To:      contactTo,
		ReplyTo: email,
		Subject: subject,
		Text:    message,
	}

	resendReqBodyBytes, err := json.Marshal(resendReqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Resend email request body: %w", err)
	}

	resendClient := &http.Client{}
	resendApiReq, err := http.NewRequest("POST", resendApiUrl, bytes.NewBuffer(resendReqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating Resend API request: %w", err)
	}
	resendApiReq.Header.Set("Content-Type", "application/json")
	resendApiReq.Header.Set("Authorization", "Bearer "+resendApiKey)

	resendApiRes, err := resendClient.Do(resendApiReq)
	if err != nil {
		return nil, fmt.Errorf("error making Resend API request to %s: %w", resendApiUrl, err)
	}
	defer resendApiRes.Body.Close()

	resendResponseBody, _ := io.ReadAll(resendApiRes.Body)
	if resendApiRes.StatusCode != http.StatusOK {
		log.Printf("Resend API returned non-200 status: %d, response: %s", resendApiRes.StatusCode, string(resendResponseBody))
		return nil, fmt.Errorf("Resend API error: status %d, response: %s", resendApiRes.StatusCode, string(resendResponseBody))
	}

	var resendApiResponse any
	if err := json.Unmarshal(resendResponseBody, &resendApiResponse); err != nil {
		log.Printf("Error unmarshalling Resend API response: %v, raw body: %s", err, string(resendResponseBody))
		return nil, fmt.Errorf("error unmarshalling Resend API response: %w", err)
	}

	return resendApiResponse, nil
}
