package config

import (
	"os"
	"strconv"
)

type Config struct {
	Domain              string
	SMTPPort           int
	CloudflareEmail    string
	CloudflareAPIToken string
	CloudflareZoneID   string
	CloudflareAccountID string
	ServerIP           string
}

func Load() *Config {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "25"))
	
	return &Config{
		Domain:              getEnv("DOMAIN", ""),
		SMTPPort:           smtpPort,
		CloudflareEmail:    getEnv("CLOUDFLARE_EMAIL", ""),
		CloudflareAPIToken: getEnv("CLOUDFLARE_API_TOKEN", ""),
		CloudflareZoneID:   getEnv("CLOUDFLARE_ZONE_ID", ""),
		CloudflareAccountID: getEnv("CLOUDFLARE_ACCOUNT_ID", ""),
		ServerIP:           getEnv("SERVER_IP", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
