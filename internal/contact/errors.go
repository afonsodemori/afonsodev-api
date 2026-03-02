package contact

import "errors"

var (
	ErrMissingFields      = errors.New("contact.form.missing_fields")
	ErrMissingToken       = errors.New("contact.form.captcha.missing")
	ErrInvalidToken       = errors.New("contact.form.captcha.invalid")
	ErrUnknownChallenger  = errors.New("contact.form.unknown_challenger")
	ErrChallengeFailed    = errors.New("contact.form.challenge_failed")
	ErrEmailSendFailed    = errors.New("contact.form.email_failed")
)
