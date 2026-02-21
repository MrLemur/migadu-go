package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// validateForwardingExpiry validates forwarding expiry date requirements.
func validateForwardingExpiry(f *Forwarding) error {
	if f.ExpiresOn != "" {
		expiryDate, err := time.Parse("2006-01-02", f.ExpiresOn)
		if err != nil {
			return fmt.Errorf("invalid expires_on date format (expected YYYY-MM-DD): %w", err)
		}

		today := time.Now().Truncate(24 * time.Hour)
		if expiryDate.Before(today) {
			return fmt.Errorf("expires_on must be a future date, got: %s", f.ExpiresOn)
		}
	}
	return nil
}

// Forwarding represents a mailbox forwarding in the Migadu API.
type Forwarding struct {
	Address            string `json:"address,omitempty" api:"create-only"`
	BlockedAt          string `json:"blocked_at,omitempty" api:"read-only"`
	ConfirmationSentAt string `json:"confirmation_sent_at,omitempty" api:"read-only"`
	ConfirmedAt        string `json:"confirmed_at,omitempty" api:"read-only"`
	ExpiresOn          string `json:"expires_on,omitempty"`
	IsActive           bool   `json:"is_active,omitempty"`
	RemoveUponExpiry   bool   `json:"remove_upon_expiry,omitempty"`
}

// ListForwardings lists all forwardings for the given domain and mailbox.
// It returns a slice of Forwarding structs and any error encountered.
func (c *Client) ListForwardings(ctx context.Context, d *Domain, mb *Mailbox) ([]Forwarding, error) {
	var forwardingList struct {
		Forwardings []Forwarding `json:"forwardings,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &forwardingList); err != nil {
		return nil, err
	}

	return forwardingList.Forwardings, nil
}

// GetForwarding retrieves a single forwarding given the domain, mailbox and a Forwarding with Address set.
// It returns a pointer to a Forwarding struct and any error encountered.
func (c *Client) GetForwarding(ctx context.Context, d *Domain, mb *Mailbox, f *Forwarding) (*Forwarding, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(f.Address)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, f); err != nil {
		return nil, err
	}

	return f, nil
}

// NewForwarding creates a new forwarding.
// It returns a pointer to a Forwarding struct and any error encountered.
func (c *Client) NewForwarding(ctx context.Context, d *Domain, mb *Mailbox, f *Forwarding) (*Forwarding, error) {
	if err := validateForwardingExpiry(f); err != nil {
		return nil, err
	}

	transformed, err := Transform(*f, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, f); err != nil {
		return nil, err
	}
	return f, nil
}

// UpdateForwarding updates a forwarding in place given the domain, mailbox and a pointer to a Forwarding struct.
// It returns a pointer to a new Forwarding struct and any error encountered.
func (c *Client) UpdateForwarding(ctx context.Context, d *Domain, mb *Mailbox, f *Forwarding) (*Forwarding, error) {
	if err := validateForwardingExpiry(f); err != nil {
		return nil, err
	}

	transformed, err := Transform(*f, "update")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(f.Address)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, f); err != nil {
		return nil, err
	}
	return f, nil
}

// DeleteForwarding deletes a forwarding given the domain, mailbox and a pointer to a Forwarding struct.
// It returns any error encountered.
func (c *Client) DeleteForwarding(ctx context.Context, d *Domain, mb *Mailbox, f *Forwarding) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/forwardings/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(f.Address)))
	if err != nil {
		return err
	}
	return nil
}
