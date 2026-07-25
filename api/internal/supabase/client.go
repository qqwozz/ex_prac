package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	anonKey    string
	serviceKey string
	http       *http.Client
}

func NewClient(baseURL, anonKey, serviceKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		anonKey:    anonKey,
		serviceKey: serviceKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) headers() map[string]string {
	return map[string]string{
		"apikey":       c.anonKey,
		"Authorization": "Bearer " + c.anonKey,
		"Content-Type":  "application/json",
	}
}

func (c *Client) headersService() map[string]string {
	return map[string]string{
		"apikey":       c.serviceKey,
		"Authorization": "Bearer " + c.serviceKey,
		"Content-Type":  "application/json",
	}
}

// RawQuery — GET-запрос к PostgREST
func (c *Client) RawQuery(endpoint string, useServiceRole bool) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Query — запрос с распаковкой JSON
func (c *Client) Query(endpoint string, useServiceRole bool, result interface{}) error {
	body, err := c.RawQuery(endpoint, useServiceRole)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

// Patch — PATCH-запрос к PostgREST (обновление записей)
func (c *Client) Patch(endpoint string, useServiceRole bool, payload interface{}) error {
	url := fmt.Sprintf("%s/rest/v1/%s", c.baseURL, endpoint)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	h := c.headers()
	if useServiceRole {
		h = c.headersService()
	}
	h["Prefer"] = "return=minimal"
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
