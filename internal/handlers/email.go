package handlers

import (
	"net/http"
	"strconv"

	"tele-temp-mail/internal/models"
	"tele-temp-mail/internal/services"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) CreateEmail(c *gin.Context) {
	var req models.CreateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Format permintaan tidak valid",
		})
		return
	}

	email, err := h.emailService.CreateEmail(req.Domain, req.TTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "creation_failed",
			Message: "Gagal membuat email sementara",
		})
		return
	}

	response := models.CreateEmailResponse{
		Email:   *email,
		Message: "Email sementara berhasil dibuat",
	}

	c.JSON(http.StatusCreated, response)
}

func (h *EmailHandler) GetEmail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_id",
			Message: "ID email diperlukan",
		})
		return
	}

	email, err := h.emailService.GetEmail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "email_not_found",
			Message: "Email tidak ditemukan atau sudah kedaluwarsa",
		})
		return
	}

	c.JSON(http.StatusOK, email)
}

func (h *EmailHandler) GetMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_id",
			Message: "ID email diperlukan",
		})
		return
	}

	// Check if email exists first
	_, err := h.emailService.GetEmail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "email_not_found",
			Message: "Email tidak ditemukan atau sudah kedaluwarsa",
		})
		return
	}

	messages, err := h.emailService.GetMessages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "fetch_failed",
			Message: "Gagal mengambil pesan",
		})
		return
	}

	// Add pagination support
	page := 1
	limit := 50
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(messages) {
		c.JSON(http.StatusOK, gin.H{
			"messages": []models.Message{},
			"total":    len(messages),
			"page":     page,
			"limit":    limit,
		})
		return
	}

	if end > len(messages) {
		end = len(messages)
	}

	paginatedMessages := messages[start:end]

	c.JSON(http.StatusOK, gin.H{
		"messages": paginatedMessages,
		"total":    len(messages),
		"page":     page,
		"limit":    limit,
	})
}

func (h *EmailHandler) DeleteEmail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_id",
			Message: "ID email diperlukan",
		})
		return
	}

	err := h.emailService.DeleteEmail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "email_not_found",
			Message: "Email tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email berhasil dihapus",
	})
}
