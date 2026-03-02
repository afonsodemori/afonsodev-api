package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SendRequest struct {
	From    string
	To      string
	ReplyTo string
	Subject string
	Text    string
}

type resendPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	ReplyTo string `json:"reply_to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type Sender interface {
	Send(ctx context.Context, req SendRequest) (any, error)
}

type ResendClient struct {
	apiKey string
}

func NewResendClient(apiKey string) *ResendClient {
	return &ResendClient{apiKey: apiKey}
}

func (c *ResendClient) Send(ctx context.Context, req SendRequest) (any, error) {
	apiURL := "https://api.resend.com/emails"

	payload := resendPayload{
		From:    req.From,
		To:      req.To,
		ReplyTo: req.ReplyTo,
		Subject: req.Subject,
		Text:    req.Text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Marshalling resend request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Creating resend request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Resend request to %s: %w", apiURL, err)
	}
	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Printf("Resend API returned non-200 status: %d, response: %s", res.StatusCode, string(responseBody))
		return nil, fmt.Errorf("Resend API error: status %d, response: %s", res.StatusCode, string(responseBody))
	}

	var response any
	if err := json.Unmarshal(responseBody, &response); err != nil {
		log.Printf("Error unmarshalling Resend API response: %v, raw body: %s", err, string(responseBody))
		return nil, fmt.Errorf("Unmarshalling resend response: %w", err)
	}

	return response, nil
}
