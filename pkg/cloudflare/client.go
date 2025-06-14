package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	email     string
	apiToken  string
	zoneID    string
	accountID string
	baseURL   string
}

type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type CloudflareResponse struct {
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
	Result  interface{} `json:"result"`
}

func NewClient(email, apiToken, zoneID, accountID string) *Client {
	return &Client{
		email:     email,
		apiToken:  apiToken,
		zoneID:    zoneID,
		accountID: accountID,
		baseURL:   "https://api.cloudflare.com/v4",
	}
}

func (c *Client) SetupDNSRecords(domain, serverIP string) error {
	// Create MX record
	mxRecord := DNSRecord{
		Type:    "MX",
		Name:    domain,
		Content: fmt.Sprintf("10 %s", domain),
		TTL:     300,
	}

	if err := c.createOrUpdateDNSRecord(mxRecord); err != nil {
		return fmt.Errorf("failed to create MX record: %v", err)
	}

	// Create A record for the domain
	aRecord := DNSRecord{
		Type:    "A",
		Name:    domain,
		Content: serverIP,
		TTL:     300,
	}

	if err := c.createOrUpdateDNSRecord(aRecord); err != nil {
		return fmt.Errorf("failed to create A record: %v", err)
	}

	// Create SPF record
	spfRecord := DNSRecord{
		Type:    "TXT",
		Name:    domain,
		Content: fmt.Sprintf("v=spf1 ip4:%s ~all", serverIP),
		TTL:     300,
	}

	if err := c.createOrUpdateDNSRecord(spfRecord); err != nil {
		return fmt.Errorf("failed to create SPF record: %v", err)
	}

	return nil
}

func (c *Client) createOrUpdateDNSRecord(record DNSRecord) error {
	// First, try to find existing record
	existingRecords, err := c.listDNSRecords(record.Name, record.Type)
	if err != nil {
		return err
	}

	if len(existingRecords) > 0 {
		// Update existing record
		record.ID = existingRecords[0].ID
		return c.updateDNSRecord(record)
	}

	// Create new record
	return c.createDNSRecord(record)
}

func (c *Client) listDNSRecords(name, recordType string) ([]DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?name=%s&type=%s", c.baseURL, c.zoneID, name, recordType)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Success bool        `json:"success"`
		Result  []DNSRecord `json:"result"`
		Errors  []string    `json:"errors"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("cloudflare API error: %v", response.Errors)
	}

	return response.Result, nil
}

func (c *Client) createDNSRecord(record DNSRecord) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.baseURL, c.zoneID)
	
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response CloudflareResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("cloudflare API error: %v", response.Errors)
	}

	return nil
}

func (c *Client) updateDNSRecord(record DNSRecord) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.baseURL, c.zoneID, record.ID)
	
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response CloudflareResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("cloudflare API error: %v", response.Errors)
	}

	return nil
}
