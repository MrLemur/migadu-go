package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ListDomains lists all domains for the account.
func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var domainList struct {
		Domains []Domain `json:"domains,omitempty"`
	}

	resp, err := c.Get(ctx, "domains")
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

	return domainList.Domains, nil
}

// GetDomain retrieves a single domain.
func (c *Client) GetDomain(ctx context.Context, domain *Domain) (*Domain, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s", domain.Name))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// NewDomain creates a new domain.
func (c *Client) NewDomain(ctx context.Context, domain *Domain) (*Domain, error) {
	jsonBody, err := json.Marshal(domain)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, "domains", jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// UpdateDomain updates a domain in place given a pointer to a Domain struct.
func (c *Client) UpdateDomain(ctx context.Context, d *Domain) (*Domain, error) {
	jsonBody, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	resp, err := c.Patch(ctx, fmt.Sprintf("domains/%s", d.Name), jsonBody)
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
func (c *Client) GetDomainRecords(ctx context.Context, domain *Domain) ([]DNSRecord, error) {
	var recordList struct {
		Records []DNSRecord `json:"records,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/records", domain.Name))
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

	return recordList.Records, nil
}

// GetDomainDiagnostics runs DNS validation checks for a domain.
func (c *Client) GetDomainDiagnostics(ctx context.Context, domain *Domain) (*DomainDiagnostics, error) {
	var diagnostics DomainDiagnostics

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/diagnostics", domain.Name))
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
func (c *Client) ActivateDomain(ctx context.Context, domain *Domain) (*Domain, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/activate", domain.Name))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, domain); err != nil {
		return nil, err
	}

	return domain, nil
}
