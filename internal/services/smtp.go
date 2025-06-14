package services

import (
	"fmt"
	"io"
	"log"
	"strings"

	"tele-temp-mail/internal/config"
	"tele-temp-mail/internal/models"

	"github.com/emersion/go-message"
	"github.com/emersion/go-smtp"
)

type SMTPService struct {
	config       *config.Config
	emailService *EmailService
	server       *smtp.Server
}

func NewSMTPService(cfg *config.Config, emailService *EmailService) *SMTPService {
	return &SMTPService{
		config:       cfg,
		emailService: emailService,
	}
}

func (s *SMTPService) Start() error {
	s.server = smtp.NewServer(&Backend{
		emailService: s.emailService,
		domain:       s.config.Domain,
	})

	s.server.Addr = fmt.Sprintf(":%d", s.config.SMTPPort)
	s.server.Domain = s.config.Domain
	s.server.AllowInsecureAuth = true

	log.Printf("Starting SMTP server on port %d", s.config.SMTPPort)
	return s.server.ListenAndServe()
}

func (s *SMTPService) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// Backend implements SMTP server backend
type Backend struct {
	emailService *EmailService
	domain       string
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{
		emailService: bkd.emailService,
		domain:       bkd.domain,
	}, nil
}

// Session implements SMTP session
type Session struct {
	emailService *EmailService
	domain       string
	from         string
	to           []string
}

func (s *Session) AuthPlain(username, password string) error {
	return nil // Accept all auth for temporary emails
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// Check if the recipient email exists in our system
	if !strings.HasSuffix(to, "@"+s.domain) {
		return fmt.Errorf("recipient not found")
	}

	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	// Parse the email message
	msg, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("failed to parse message: %v", err)
	}

	// Extract headers
	headers := make(map[string][]string)
	for field := msg.Header.Fields(); field.Next(); {
		key := field.Key()
		value := field.Value()
		headers[key] = append(headers[key], value)
	}

	// Get subject
	subject := ""
	if subjectHeader := msg.Header.Get("Subject"); subjectHeader != "" {
		subject = subjectHeader
	}

	// Read body
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to read message body: %v", err)
	}

	// Create message for each recipient
	for _, recipient := range s.to {
		message := &models.Message{
			From:    s.from,
			To:      recipient,
			Subject: subject,
			Body:    string(body),
			Headers: headers,
		}

		// Try to add message to email
		if err := s.emailService.AddMessage(recipient, message); err != nil {
			log.Printf("Failed to add message to %s: %v", recipient, err)
			// Don't return error to avoid bouncing legitimate emails
		} else {
			log.Printf("Message received for %s from %s: %s", recipient, s.from, subject)
		}
	}

	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *Session) Logout() error {
	return nil
}
