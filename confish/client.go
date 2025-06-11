package confish

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ConfishConfig struct {
	URL         string
	AppID       string
	AppSecret   string
	WebhookPath string
}

// Client represents a confish client for configuration and logging
type Client struct {
	cfg *ConfishConfig
}

// LogLevel represents the logging level
type LogLevel string

const (
	// LogLevelDebug is for detailed debug information
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo is for general information
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn is for warning messages
	LogLevelWarn LogLevel = "warn"
	// LogLevelError is for error messages
	LogLevelError LogLevel = "error"
	// LogLevelCritical is for critical error messages
	LogLevelCritical LogLevel = "critical"
)

// LogPayload represents the payload for the logging endpoint
type LogPayload struct {
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
}

// NewClient creates a new Confish client
func NewClient(cfg *ConfishConfig) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	if cfg.URL == "" {
		return nil, errors.New("config.URL cannot be empty")
	}

	if cfg.AppID == "" {
		return nil, errors.New("config.AppID cannot be empty")
	}

	if cfg.AppSecret == "" {
		return nil, errors.New("config.AppSecret cannot be empty")
	}

	return &Client{cfg: cfg}, nil
}

// GetConfig retrieves a configuration from the Confish API and unmarshals it into the provided type
func (c *Client) GetConfig(configID string, result interface{}) error {
	url := fmt.Sprintf("%s/c/%s", c.cfg.URL, configID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("App-ID", c.cfg.AppID)
	req.Header.Add("App-Secret", c.cfg.AppSecret)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Log sends a log message to the Confish logging endpoint
func (c *Client) Log(level LogLevel, message string) error {
	return c.LogWithURL(level, message)
}

// LogWithURL sends a log message to a specific Confish logging endpoint URL
func (c *Client) LogWithURL(level LogLevel, message string) error {
	payload := LogPayload{
		Level:   level,
		Message: message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal log payload: %w", err)
	}

	url := fmt.Sprintf("%s/a/%s/log", c.cfg.URL, c.cfg.AppID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create log request: %w", err)
	}

	// Add headers
	req.Header.Add("App-ID", c.cfg.AppID)
	req.Header.Add("App-Secret", c.cfg.AppSecret)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send log: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response for log: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// WebhookPayload represents a webhook payload type received from confish
type WebhookPayload struct {
	Event         string              `json:"event"`
	Configuration ConfigurationObject `json:"configuration"`
}

// ConfigurationObject represents a configuration object received from confish
type ConfigurationObject struct {
	Name   string          `json:"name"`
	Values json.RawMessage `json:"values"`
}

// ProcessWebhookPayload processes a webhook payload and returns the configuration values
func (c *Client) ProcessWebhookPayload(payload WebhookPayload, result interface{}) error {
	// Check if this is a configuration update event
	if payload.Event != "configuration.updated" {
		return fmt.Errorf("unsupported event type: %s", payload.Event)
	}

	if err := json.Unmarshal(payload.Configuration.Values, result); err != nil {
		return fmt.Errorf("failed to unmarshal configuration values: %w", err)
	}

	return nil
}

// Debug logs a debug message
func (c *Client) Debug(message string) error {
	return c.Log(LogLevelDebug, message)
}

// Info logs an info message
func (c *Client) Info(message string) error {
	return c.Log(LogLevelInfo, message)
}

// Warn logs a warning message
func (c *Client) Warn(message string) error {
	return c.Log(LogLevelWarn, message)
}

// Error logs an error message
func (c *Client) Error(message string) error {
	return c.Log(LogLevelError, message)
}

// Critical logs a critical message
func (c *Client) Critical(message string) error {
	return c.Log(LogLevelCritical, message)
}
