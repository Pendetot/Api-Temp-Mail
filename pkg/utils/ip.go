package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetPublicIP retrieves the public IP address of the server
func GetPublicIP() (string, error) {
	// List of IP detection services
	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://ipecho.net/plain",
		"https://myexternalip.com/raw",
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, service := range services {
		ip, err := getIPFromService(client, service)
		if err == nil && ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("failed to get public IP from all services")
}

func getIPFromService(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty response")
	}

	return ip, nil
}
