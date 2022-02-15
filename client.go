package migadu

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

// httPClient implements the most basic function of http.Client.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a client for working with Migadu API.
type Client struct {
	Email      string
	APIKey     string
	Domain     string
	Timeout    int
	BaseURL    string
	HTTPClient httpClient
}

// testAuth tests that credentials are valid before creating a client.
// It returns any error encountered.
func (c *Client) testAuth() error {
	ctx := context.Background()
	_, err := c.Get(ctx, "mailboxes")
	if err != nil {
		return fmt.Errorf("Authorisation not valid. Error: %w", err)
	}
	return nil
}

// New creates a new client one domain on Migadu given the admin email and API key.
// It returns a pointer to the new client and any error encountered.
func New(email string, apiKey string, domain string) (*Client, error) {
	baseURL := fmt.Sprintf("https://api.migadu.com/v1/domains/%s", domain)

	client := &Client{
		Email:      email,
		APIKey:     apiKey,
		Domain:     domain,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}

	err := client.testAuth()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Get is a convenience function to GET an HTTP resource given a path.
// It can be used to access the API directly if the existing methods are not enough.
// It returns a pointer to an http.Response and any error encountered.
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*120))
	defer cancel()

	urlPath := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	req.SetBasicAuth(c.Email, c.APIKey)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d is not 200", resp.StatusCode)
	}

	return resp, nil
}

// Post is a convenience function to POST to an HTTP resource given a path and a byte slice.
// It can be used to access the API directly if the existing methods are not enough.
// It returns a pointer to an http.Response and any error encountered.
func (c *Client) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*120))
	defer cancel()

	urlPath := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(body))
	req.SetBasicAuth(c.Email, c.APIKey)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: status code %d is not 200", resp.StatusCode)
	}

	return resp, nil
}

// Put is a convenience function to PUT to an HTTP resource given a path and a byte slice.
// It can be used to access the API directly if the existing methods are not enough.
// It returns a pointer to an http.Response and any error encountered.
func (c *Client) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*120))
	defer cancel()

	urlPath := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, urlPath, bytes.NewBuffer(body))
	req.SetBasicAuth(c.Email, c.APIKey)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: status code %d is not 200", resp.StatusCode)
	}

	return resp, nil
}

// Delete is a convenience function to DELETE an HTTP resource given a path.
// It can be used to access the API directly if the existing methods are not enough.
// It returns a pointer to an http.Response and any error encountered.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*120))
	defer cancel()

	urlPath := fmt.Sprintf("%s/%s", c.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, urlPath, nil)
	req.SetBasicAuth(c.Email, c.APIKey)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: status code %d is not 200", resp.StatusCode)
	}

	return resp, nil
}
