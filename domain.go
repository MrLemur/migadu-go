package migadu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Domain represents a domain in the Migadu API.
type Domain struct {
	Name                 string   `json:"name,omitempty"`
	State                string   `json:"state,omitempty"`
	Description          string   `json:"description,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	SpamAggressiveness   string   `json:"spam_aggressiveness,omitempty"`
	GreylistingEnabled   bool     `json:"greylisting_enabled,omitempty"`
	MXProxyEnabled       bool     `json:"mx_proxy_enabled,omitempty"`
	HostedDNS            bool     `json:"hosted_dns,omitempty"`
	SenderDenylist       []string `json:"sender_denylist,omitempty"`
	RecipientDenylist    []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations []string `json:"catchall_destinations,omitempty"`
}

// DNSRecord represents a DNS record needed for domain configuration.
type DNSRecord struct {
	Type     string `json:"type,omitempty"`
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	Priority int    `json:"priority,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
}

// DomainDiagnostics represents DNS validation results for a domain.
type DomainDiagnostics struct {
	MX    []string `json:"mx,omitempty"`
	SPF   []string `json:"spf,omitempty"`
	DKIM  []string `json:"dkim,omitempty"`
	DMARC []string `json:"dmarc,omitempty"`
}

// apiRequest makes an HTTP request to the root API endpoint (not domain-scoped).
func (c *Client) apiRequest(ctx context.Context, method string, path string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*120))
	defer cancel()

	urlPath := fmt.Sprintf("https://api.migadu.com/v1/%s", path)
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequestWithContext(ctx, method, urlPath, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, urlPath, nil)
	}

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Email, c.APIKey)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
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

// ListDomains lists all domains for the account.
func (c *Client) ListDomains(ctx context.Context) (*[]Domain, error) {
	var domainList struct {
		Domains []Domain `json:"domains,omitempty"`
	}

	resp, err := c.apiRequest(ctx, http.MethodGet, "domains", nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &domainList); err != nil {
		return nil, err
	}

	return &domainList.Domains, nil
}

// GetDomain retrieves a single domain given its name.
func (c *Client) GetDomain(ctx context.Context, name string) (*Domain, error) {
	var domain Domain

	resp, err := c.apiRequest(ctx, http.MethodGet, fmt.Sprintf("domains/%s", name), nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// NewDomain creates a new domain given its name.
func (c *Client) NewDomain(ctx context.Context, name string) (*Domain, error) {
	var domain = Domain{Name: name}

	jsonBody, err := json.Marshal(domain)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiRequest(ctx, http.MethodPost, "domains", jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// UpdateDomain updates a domain in place given a pointer to a Domain struct.
func (c *Client) UpdateDomain(ctx context.Context, d *Domain) (*Domain, error) {
	jsonBody, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiRequest(ctx, http.MethodPatch, fmt.Sprintf("domains/%s", d.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomainRecords retrieves the DNS records needed for a domain.
func (c *Client) GetDomainRecords(ctx context.Context, name string) (*[]DNSRecord, error) {
	var recordList struct {
		Records []DNSRecord `json:"records,omitempty"`
	}

	resp, err := c.apiRequest(ctx, http.MethodGet, fmt.Sprintf("domains/%s/records", name), nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &recordList); err != nil {
		return nil, err
	}

	return &recordList.Records, nil
}

// GetDomainDiagnostics runs DNS validation checks for a domain.
func (c *Client) GetDomainDiagnostics(ctx context.Context, name string) (*DomainDiagnostics, error) {
	var diagnostics DomainDiagnostics

	resp, err := c.apiRequest(ctx, http.MethodGet, fmt.Sprintf("domains/%s/diagnostics", name), nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &diagnostics); err != nil {
		return nil, err
	}

	return &diagnostics, nil
}

// ActivateDomain activates a domain after DNS setup is complete.
func (c *Client) ActivateDomain(ctx context.Context, name string) (*Domain, error) {
	var domain Domain

	resp, err := c.apiRequest(ctx, http.MethodGet, fmt.Sprintf("domains/%s/activate", name), nil)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}
