package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// Forwarding represents a mailbox forwarding in the Migadu API.
type Forwarding struct {
	Address            string  `json:"address,omitempty"`
	BlockedAt          *string `json:"blocked_at,omitempty"`
	ConfirmationSentAt *string `json:"confirmation_sent_at,omitempty"`
	ConfirmedAt        *string `json:"confirmed_at,omitempty"`
	ExpiresOn          *string `json:"expires_on,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
	RemoveUponExpiry   *bool   `json:"remove_upon_expiry,omitempty"`
}

// ListForwardings lists all forwardings for the given domain and mailbox.
// It returns a slice of Forwarding structs and any error encountered.
func (c *Client) ListForwardings(ctx context.Context, domain *Domain, mailbox string) ([]Forwarding, error) {

	var forwardingList struct {
		Forwardings []Forwarding `json:"forwardings,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", domain.Name, mailbox))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &forwardingList); err != nil {
		return nil, err
	}

	return forwardingList.Forwardings, nil
}

// GetForwarding retrieves a single forwarding given the domain, mailbox and address.
// It returns a pointer to a Forwarding struct and any error encountered.
func (c *Client) GetForwarding(ctx context.Context, domain *Domain, mailbox string, address string) (*Forwarding, error) {

	var forwarding Forwarding

	escapedAddress := url.PathEscape(address)
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", domain.Name, mailbox, escapedAddress))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &forwarding); err != nil {
		return nil, err
	}

	return &forwarding, nil
}

// NewForwarding creates a new forwarding.
// It returns a pointer to a Forwarding struct and any error encountered.
func (c *Client) NewForwarding(ctx context.Context, domain *Domain, mailbox string, forwarding *Forwarding) (*Forwarding, error) {
	jsonBody, err := json.Marshal(forwarding)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", domain.Name, mailbox), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, forwarding); err != nil {
		return nil, err
	}
	return forwarding, nil
}

// UpdateForwarding updates a forwarding in place given the domain, mailbox and a pointer to a Forwarding struct.
// It returns a pointer to a new Forwarding struct and any error encountered.
func (c *Client) UpdateForwarding(ctx context.Context, domain *Domain, mailbox string, f *Forwarding) (*Forwarding, error) {
	jsonBody, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	escapedAddress := url.PathEscape(f.Address)
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", domain.Name, mailbox, escapedAddress), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &f); err != nil {
		return nil, err
	}
	return f, nil
}

// DeleteForwarding deletes a forwarding given the domain, mailbox and a pointer to a Forwarding struct.
// It returns any error encountered.
func (c *Client) DeleteForwarding(ctx context.Context, domain *Domain, mailbox string, f *Forwarding) error {
	escapedAddress := url.PathEscape(f.Address)
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", domain.Name, mailbox, escapedAddress))
	if err != nil {
		return err
	}
	return nil
}
