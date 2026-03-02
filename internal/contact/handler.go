package contact

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleSendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, SendEmailResponse{
			Success: false,
			Message: "Method Not Allowed",
		})
		return
	}

	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Error:   "Invalid request payload",
		})
		return
	}

	result, err := h.service.SendEmail(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, SendEmailResponse{
		Success: true,
		Message: "contact.form.success",
		Details: result.APIResponse,
	})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrMissingFields):
		writeJSON(w, http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Message: err.Error(),
		})
	case errors.Is(err, ErrMissingToken):
		writeJSON(w, http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		})
	case errors.Is(err, ErrUnknownChallenger):
		writeJSON(w, http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		})
	case errors.Is(err, ErrInvalidToken):
		details := extractDetails(err)
		writeJSON(w, http.StatusBadRequest, SendEmailResponse{
			Success: false,
			Error:   ErrInvalidToken.Error(),
			Details: details,
		})
	case errors.Is(err, ErrChallengeFailed):
		writeJSON(w, http.StatusInternalServerError, SendEmailResponse{
			Success: false,
			Error:   "Challenge verification failed",
		})
	case errors.Is(err, ErrEmailSendFailed):
		writeJSON(w, http.StatusInternalServerError, SendEmailResponse{
			Success: false,
			Error:   "Failed to send email",
			Details: extractDetails(err),
		})
	default:
		writeJSON(w, http.StatusInternalServerError, SendEmailResponse{
			Success: false,
			Error:   "Internal server error",
		})
	}
}

func extractDetails(err error) string {
	msg := err.Error()
	parts := strings.SplitN(msg, ": ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return msg
}

// TODO: Maybe a helper? /app/internal/pkg/response/json.go
// Have in mind: Go-style -> a little copying is better than a little dependency
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
