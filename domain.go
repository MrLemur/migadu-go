package migadu

import (
	"context"
	"encoding/json"
	"fmt"
)

// validateDomainCreate validates domain creation parameters.
func validateDomainCreate(d *Domain) error {
	if d.HostedDNS {
		return fmt.Errorf("hosted_dns: true is not recommended as Migadu plans to discontinue this service")
	}
	return nil
}

// Domain represents a domain in the Migadu API.
//
// SpamAggressiveness is a string with the following valid values:
//
//	"paranoid"    (most aggressive filtering)
//	"aggressive"
//	"default"
//	"suspicious"
//	"permissive"  (least aggressive filtering)
type Domain struct {
	Name                             string   `json:"name,omitempty" api:"create-only"`
	ActivatedAt                      string   `json:"activated_at,omitempty" api:"read-only"`
	DeactivatedAt                    string   `json:"deactivated_at,omitempty" api:"read-only"`
	State                            string   `json:"state,omitempty" api:"read-only"`
	Description                      string   `json:"description,omitempty"`
	Tags                             []string `json:"tags,omitempty"`
	CanSend                          bool     `json:"can_send,omitempty" api:"read-only"`
	CanReceive                       bool     `json:"can_receive,omitempty" api:"read-only"`
	CanAccess                        bool     `json:"can_access,omitempty" api:"read-only"`
	SpamAggressiveness               string   `json:"spam_aggressiveness,omitempty"`
	GreylistingEnabled               bool     `json:"greylisting_enabled,omitempty"`
	JunkSubjectKeywordSpam           bool     `json:"junk_subject_keyword_spam,omitempty"`
	SubjectRewritingEnabled          bool     `json:"subject_rewriting_enabled,omitempty"`
	MXProxyEnabled                   bool     `json:"mx_proxy_enabled,omitempty"`
	HostedDNS                        bool     `json:"hosted_dns,omitempty"`
	MailboxDefaultImapEnabled        bool     `json:"mailbox_default_imap_enabled,omitempty"`
	MailboxDefaultIncomingLimit      int      `json:"mailbox_default_incoming_limit,omitempty"`
	MailboxDefaultManagesieveEnabled bool     `json:"mailbox_default_managesieve_enabled,omitempty"`
	MailboxDefaultOutgoingLimit      int      `json:"mailbox_default_outgoing_limit,omitempty"`
	MailboxDefaultPop3Enabled        bool     `json:"mailbox_default_pop3_enabled,omitempty"`
	MailboxDefaultReceivingEnabled   bool     `json:"mailbox_default_receiving_enabled,omitempty"`
	MailboxDefaultSendingEnabled     bool     `json:"mailbox_default_sending_enabled,omitempty"`
	MailboxDefaultStorageLimit       int      `json:"mailbox_default_storage_limit,omitempty"`
	SenderAllowlist                  []string `json:"sender_allowlist,omitempty"`
	SenderDenylist                   []string `json:"sender_denylist,omitempty"`
	RecipientDenylist                []string `json:"recipient_denylist,omitempty"`
	CatchallDestinations             []string `json:"catchall_destinations,omitempty"`
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
	if err = readAndUnmarshal(resp.Body, &domainList); err != nil {
		return nil, err
	}

	return domainList.Domains, nil
}

// GetDomain retrieves a single domain.
func (c *Client) GetDomain(ctx context.Context, d *Domain) (*Domain, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s", d.Name))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, d); err != nil {
		return nil, err
	}

	return d, nil
}

// NewDomain creates a new domain.
func (c *Client) NewDomain(ctx context.Context, d *Domain) (*Domain, error) {
	if err := validateDomainCreate(d); err != nil {
		return nil, err
	}

	transformed, err := Transform(*d, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "domains", jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, d); err != nil {
		return nil, err
	}

	return d, nil
}

// UpdateDomain updates a domain in place given a pointer to a Domain struct.
func (c *Client) UpdateDomain(ctx context.Context, d *Domain) (*Domain, error) {
	transformed, err := Transform(*d, "update")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Patch(ctx, fmt.Sprintf("domains/%s", escapePathSegment(d.Name)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomainRecords retrieves the DNS records needed for a domain.
func (c *Client) GetDomainRecords(ctx context.Context, d *Domain) ([]DNSRecord, error) {
	var apiResponse struct {
		DomainName      string      `json:"domain_name"`
		SPF             *DNSRecord  `json:"spf"`
		DKIM            []DNSRecord `json:"dkim"`
		DMARC           *DNSRecord  `json:"dmarc"`
		DNSVerification *DNSRecord  `json:"dns_verification"`
		MXRecords       []DNSRecord `json:"mx_records"`
		AutoConfig      *DNSRecord  `json:"autoconfig"`
		AutoDiscover    *DNSRecord  `json:"autodiscover"`
		SRV             []DNSRecord `json:"srv"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/records", d.Name))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &apiResponse); err != nil {
		return nil, err
	}

	var records []DNSRecord
	if apiResponse.SPF != nil {
		records = append(records, *apiResponse.SPF)
	}
	records = append(records, apiResponse.DKIM...)
	if apiResponse.DMARC != nil {
		records = append(records, *apiResponse.DMARC)
	}
	if apiResponse.DNSVerification != nil {
		verificationRecord := *apiResponse.DNSVerification
		records = append(records, verificationRecord)
	}
	records = append(records, apiResponse.MXRecords...)
	if apiResponse.AutoConfig != nil {
		records = append(records, *apiResponse.AutoConfig)
	}
	if apiResponse.AutoDiscover != nil {
		records = append(records, *apiResponse.AutoDiscover)
	}
	records = append(records, apiResponse.SRV...)

	return records, nil
}

// GetDomainDiagnostics runs DNS validation checks for a domain.
func (c *Client) GetDomainDiagnostics(ctx context.Context, d *Domain) (*DomainDiagnostics, error) {
	var diagnostics DomainDiagnostics

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/diagnostics", d.Name))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &diagnostics); err != nil {
		return nil, err
	}

	return &diagnostics, nil
}

// ActivateDomain activates a domain after DNS setup is complete.
func (c *Client) ActivateDomain(ctx context.Context, d *Domain) (*Domain, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/activate", d.Name))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, d); err != nil {
		return nil, err
	}

	return d, nil
}
