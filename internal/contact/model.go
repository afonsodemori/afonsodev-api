package contact

type SendEmailRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	Challenger string `json:"challenger"`
}

type SendEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Details any    `json:"details,omitempty"`
}
