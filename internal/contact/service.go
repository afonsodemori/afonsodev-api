package contact

import (
	"context"
	"fmt"
	"log"

	"github.com/afonsodemori/afonsodev-api/internal/challenge"
	"github.com/afonsodemori/afonsodev-api/internal/email"
)

type ChallengeVerifier interface {
	Verify(ctx context.Context, token string) (*challenge.Result, error)
}

type EmailSender interface {
	Send(ctx context.Context, req email.SendRequest) (any, error)
}

type ServiceConfig struct {
	IsDevelopment bool
	ContactFrom   string
	ContactTo     string
}

type Service struct {
	config    ServiceConfig
	recaptcha ChallengeVerifier
	turnstile ChallengeVerifier
	email     EmailSender
}

func NewService(cfg ServiceConfig, recaptcha, turnstile ChallengeVerifier, emailSender EmailSender) *Service {
	return &Service{
		config:    cfg,
		recaptcha: recaptcha,
		turnstile: turnstile,
		email:     emailSender,
	}
}

type SendEmailResult struct {
	APIResponse any
}

func (s *Service) SendEmail(ctx context.Context, req SendEmailRequest) (*SendEmailResult, error) {
	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		return nil, ErrMissingFields
	}

	if req.Token == "" {
		return nil, ErrMissingToken
	}

	var verifier ChallengeVerifier
	switch req.Challenger {
	case "", "captcha":
		verifier = s.recaptcha
	case "turnstile":
		verifier = s.turnstile
	default:
		return nil, ErrUnknownChallenger
	}

	result, err := verifier.Verify(ctx, req.Token)
	if err != nil {
		log.Printf("Challenge verification error: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrChallengeFailed, err)
	}

	if !result.Success {
		log.Printf("Challenge verification unsuccessful: %v", result.Errors)
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, result.Errors)
	}

	if s.config.IsDevelopment {
		log.Printf("Development mode. Email was not sent")
		return &SendEmailResult{}, nil
	}

	apiResponse, err := s.email.Send(ctx, email.SendRequest{
		From:    fmt.Sprintf("%s - via afonso.dev <%s>", req.Name, s.config.ContactFrom),
		To:      s.config.ContactTo,
		ReplyTo: req.Email,
		Subject: req.Subject,
		Text:    req.Message,
	})
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	return &SendEmailResult{APIResponse: apiResponse}, nil
}
