package models

import (
	"time"
)

type Email struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	ID          string              `json:"id"`
	EmailID     string              `json:"email_id"`
	From        string              `json:"from"`
	To          string              `json:"to"`
	Subject     string              `json:"subject"`
	Body        string              `json:"body"`
	HTML        string              `json:"html,omitempty"`
	Attachments []Attachment        `json:"attachments,omitempty"`
	Headers     map[string][]string `json:"headers,omitempty"`
	ReceivedAt  time.Time           `json:"received_at"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Data        []byte `json:"data,omitempty"`
}

type CreateEmailRequest struct {
	Domain string `json:"domain,omitempty"`
	TTL    int    `json:"ttl,omitempty"` // Time to live in minutes, default 60
}

type CreateEmailResponse struct {
	Email   Email  `json:"email"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
