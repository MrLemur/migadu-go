package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Domain represents a domain in the Migadu API.
type Domain struct {
	Name                    string   `json:"name,omitempty"`
	State                   string   `json:"state,omitempty"`
	Description             string   `json:"description,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	CanSend                 bool     `json:"can_send,omitempty"`
	CanReceive              bool     `json:"can_receive,omitempty"`
	CanAccess               bool     `json:"can_access,omitempty"`
	SpamAggressiveness      string   `json:"spam_aggressiveness,omitempty"`
	GreylistingEnabled      bool     `json:"greylisting_enabled,omitempty"`
	JunkSubjectKeywordSpam  bool     `json:"junk_subject_keyword_spam,omitempty"`
	SubjectRewritingEnabled bool     `json:"subject_rewriting_enabled,omitempty"`
	MXProxyEnabled          bool     `json:"mx_proxy_enabled,omitempty"`
	HostedDNS               bool     `json:"hosted_dns,omitempty"`
	SenderAllowlist         []string `json:"sender_allowlist,omitempty"`
	SenderDenylist          []string `json:"sender_denylist,omitempty"`
	RecipientDenylist       []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations    []string `json:"catchall_destinations,omitempty"`
}

func (d *Domain) UnmarshalJSON(data []byte) error {
	type DomainAlias Domain
	aux := &struct {
		SpamAggressiveness   json.RawMessage `json:"spam_aggressiveness"`
		Tags                 json.RawMessage `json:"tags"`
		SenderAllowlist      json.RawMessage `json:"sender_allowlist"`
		SenderDenylist       json.RawMessage `json:"sender_denylist"`
		RecipientDenylist    json.RawMessage `json:"recipient_denylist"`
		CatchallDestinations json.RawMessage `json:"catchall_destinations"`
		*DomainAlias
	}{
		DomainAlias: (*DomainAlias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse spam_aggressiveness (can be string or number)
	if len(aux.SpamAggressiveness) > 0 && string(aux.SpamAggressiveness) != "null" {
		var str string
		if err := json.Unmarshal(aux.SpamAggressiveness, &str); err == nil {
			d.SpamAggressiveness = str
		} else {
			var num interface{}
			if err := json.Unmarshal(aux.SpamAggressiveness, &num); err == nil {
				d.SpamAggressiveness = fmt.Sprintf("%v", num)
			}
		}
	}

	// Parse list fields that can be string or array
	if tags, err := parseStringOrArray(aux.Tags); err != nil {
		return err
	} else {
		d.Tags = tags
	}

	if allowlist, err := parseStringOrArray(aux.SenderAllowlist); err != nil {
		return err
	} else {
		d.SenderAllowlist = allowlist
	}

	if denylist, err := parseStringOrArray(aux.SenderDenylist); err != nil {
		return err
	} else {
		d.SenderDenylist = denylist
	}

	if recipientDenylist, err := parseStringOrArray(aux.RecipientDenylist); err != nil {
		return err
	} else {
		d.RecipientDenylist = recipientDenylist
	}

	if catchall, err := parseStringOrArray(aux.CatchallDestinations); err != nil {
		return err
	} else {
		d.CatchallDestinations = catchall
	}

	return nil
}

func parseStringOrArray(data json.RawMessage) ([]string, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	// Try as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s == "" {
			return nil, nil
		}
		return strings.Split(s, ","), nil
	}

	// Try as array
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}

	return nil, fmt.Errorf("unsupported list format: %s", string(data))
}

// Create returns a request-safe copy for domain create operations.
func (d Domain) Create() Domain {
	d.State = ""
	d.CanSend = false
	d.CanReceive = false
	d.CanAccess = false
	return d
}

// Update returns a request-safe copy for domain update operations.
func (d Domain) Update() Domain {
	d.Name = ""
	d.State = ""
	d.CanSend = false
	d.CanReceive = false
	d.CanAccess = false
	d.JunkSubjectKeywordSpam = false
	return d
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
	domainName := d.Name
	jsonBody, err := json.Marshal(d.Update())
	if err != nil {
		return nil, err
	}
	resp, err := c.Patch(ctx, fmt.Sprintf("domains/%s", escapePathSegment(domainName)), jsonBody)
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
