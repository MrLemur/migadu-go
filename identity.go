package migadu

import (
	"context"
	"encoding/json"
	"fmt"
)

// Identity represents an identity in the Migadu API.
type Identity struct {
	Address              string `json:"address,omitempty" api:"read-only"`
	DomainName           string `json:"domain_name,omitempty" api:"read-only"`
	LocalPart            string `json:"local_part,omitempty" api:"create-only"`
	MayAccessImap        bool   `json:"may_access_imap,omitempty"`
	MayAccessManagesieve bool   `json:"may_access_managesieve,omitempty"`
	MayAccessPop3        bool   `json:"may_access_pop3,omitempty"`
	MayReceive           bool   `json:"may_receive,omitempty"`
	MaySend              bool   `json:"may_send,omitempty"`
	Name                 string `json:"name,omitempty"`
	Password             string `json:"password,omitempty"`
	PasswordUse          string `json:"password_use,omitempty" api:"read-only"`
}

// ListIdentities lists all the identities for the given domain and mailbox.
// It returns a slice of Identity structs and any error encountered.
func (c *Client) ListIdentities(ctx context.Context, d *Domain, mb *Mailbox) ([]Identity, error) {
	var identityList struct {
		Identities []Identity `json:"identities,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &identityList); err != nil {
		return nil, err
	}

	return identityList.Identities, nil
}

// GetIdentity retrieves a single identity given the domain, mailbox and an Identity with LocalPart set.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, d *Domain, mb *Mailbox, i *Identity) (*Identity, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(i.LocalPart)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, i); err != nil {
		return nil, err
	}

	return i, nil
}

// NewIdentity creates a new identity.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) NewIdentity(ctx context.Context, d *Domain, mb *Mailbox, i *Identity) (*Identity, error) {
	transformed, err := Transform(*i, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, i); err != nil {
		return nil, err
	}
	return i, nil
}

// UpdateIdentity updates an identity in place given the domain, mailbox and a pointer to an Identity struct.
// It returns a pointer to a new Identity struct and any error encountered.
func (c *Client) UpdateIdentity(ctx context.Context, d *Domain, mb *Mailbox, i *Identity) (*Identity, error) {
	transformed, err := Transform(*i, "update")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(i.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, i); err != nil {
		return nil, err
	}
	return i, nil
}

// DeleteIdentity deletes an identity given the domain, mailbox and a pointer to an Identity struct.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, d *Domain, mb *Mailbox, i *Identity) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(d.Name), escapePathSegment(mb.LocalPart), escapePathSegment(i.LocalPart)))
	if err != nil {
		return err
	}
	return nil
}
