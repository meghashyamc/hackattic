// Package httpclient provides a reusable HTTP client with configured timeouts
// and connection pooling for making HTTP requests across the application.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	defaultClient *Client
	once          sync.Once
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
	mu         sync.RWMutex
}

type Config struct {
	BaseURL             string
	Timeout             time.Duration
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	DisableKeepAlives   bool
	Headers             map[string]string
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL:             os.Getenv("BASE_URL"),
		Timeout:             30 * time.Second,
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}
}

type RequestOption func(*requestConfig)

type requestConfig struct {
	headers map[string]string
	timeout time.Duration
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(rc *requestConfig) {
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}
		maps.Copy(rc.headers, headers)
	}
}

func WithTimeout(timeout time.Duration) RequestOption {
	return func(rc *requestConfig) {
		rc.timeout = timeout
	}
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func (r *Response) JSON(v any) error {
	return json.Unmarshal(r.Body, v)
}

func Initialize(config *Config) {
	once.Do(func() {
		if config == nil {
			config = DefaultConfig()
		}

		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        config.MaxIdleConns,
			MaxConnsPerHost:     config.MaxConnsPerHost,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   config.DisableKeepAlives,
		}

		defaultClient = &Client{
			httpClient: &http.Client{
				Transport: transport,
				Timeout:   config.Timeout,
			},
			baseURL: config.BaseURL,
			headers: config.Headers,
		}
	})
}

func GetDefaultClient() *Client {
	if defaultClient == nil {
		Initialize(nil)
	}
	return defaultClient
}

func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        config.MaxIdleConns,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   config.DisableKeepAlives,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		},
		baseURL: config.BaseURL,
		headers: config.Headers,
	}
}

func (c *Client) Get(ctx context.Context, endpoint string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, endpoint, nil, opts...)
}

func (c *Client) Put(ctx context.Context, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPut, endpoint, body, opts...)
}

func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, endpoint, body, opts...)
}

func (c *Client) Delete(ctx context.Context, endpoint string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil, opts...)
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, opts ...RequestOption) (*Response, error) {
	rc := &requestConfig{
		headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt(rc)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	fullURL := c.baseURL + endpoint

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.mu.RLock()
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	c.mu.RUnlock()

	for k, v := range rc.headers {
		req.Header.Set(k, v)
	}

	if rc.timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, rc.timeout)
		defer cancel()
		req = req.WithContext(timeoutCtx)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return response, fmt.Errorf("HTTP error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return response, nil
}

// SetDefaultHeaders updates the default headers for the client
func (c *Client) SetDefaultHeaders(headers map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.headers = headers
}
