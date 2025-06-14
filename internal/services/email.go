package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"tele-temp-mail/internal/config"
	"tele-temp-mail/internal/models"

	"github.com/google/uuid"
)

type EmailService struct {
	config  *config.Config
	emails  map[string]*models.Email
	mutex   sync.RWMutex
	cleanup chan string
}

func NewEmailService(cfg *config.Config) *EmailService {
	service := &EmailService{
		config:  cfg,
		emails:  make(map[string]*models.Email),
		cleanup: make(chan string, 100),
	}

	// Start cleanup goroutine
	go service.cleanupExpiredEmails()

	return service
}

func (s *EmailService) CreateEmail(domain string, ttl int) (*models.Email, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if domain == "" {
		domain = s.config.Domain
	}

	if ttl <= 0 {
		ttl = 60 // Default 60 minutes
	}

	// Generate random email address
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random email: %v", err)
	}

	username := hex.EncodeToString(randomBytes)
	address := fmt.Sprintf("%s@%s", username, domain)

	email := &models.Email{
		ID:        uuid.New().String(),
		Address:   address,
		Domain:    domain,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Minute),
		Messages:  []models.Message{},
	}

	s.emails[email.ID] = email

	// Schedule cleanup
	go func() {
		time.Sleep(time.Duration(ttl) * time.Minute)
		s.cleanup <- email.ID
	}()

	return email, nil
}

func (s *EmailService) GetEmail(id string) (*models.Email, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	email, exists := s.emails[id]
	if !exists {
		return nil, fmt.Errorf("email not found")
	}

	if time.Now().After(email.ExpiresAt) {
		return nil, fmt.Errorf("email expired")
	}

	return email, nil
}

func (s *EmailService) GetEmailByAddress(address string) (*models.Email, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, email := range s.emails {
		if email.Address == address && time.Now().Before(email.ExpiresAt) {
			return email, nil
		}
	}

	return nil, fmt.Errorf("email not found")
}

func (s *EmailService) DeleteEmail(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.emails[id]; !exists {
		return fmt.Errorf("email not found")
	}

	delete(s.emails, id)
	return nil
}

func (s *EmailService) AddMessage(emailAddress string, message *models.Message) error {
	email, err := s.GetEmailByAddress(emailAddress)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	message.ID = uuid.New().String()
	message.EmailID = email.ID
	message.ReceivedAt = time.Now()

	email.Messages = append(email.Messages, *message)
	return nil
}

func (s *EmailService) GetMessages(emailID string) ([]models.Message, error) {
	email, err := s.GetEmail(emailID)
	if err != nil {
		return nil, err
	}

	return email.Messages, nil
}

func (s *EmailService) cleanupExpiredEmails() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case emailID := <-s.cleanup:
			s.mutex.Lock()
			delete(s.emails, emailID)
			s.mutex.Unlock()

		case <-ticker.C:
			s.mutex.Lock()
			now := time.Now()
			for id, email := range s.emails {
				if now.After(email.ExpiresAt) {
					delete(s.emails, id)
				}
			}
			s.mutex.Unlock()
		}
	}
}
