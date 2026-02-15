package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Identity represents an identity in the Migadu API.
type Identity struct {
	Address              string `json:"address,omitempty"`
	DomainName           string `json:"domain_name,omitempty"`
	FooterActive         bool   `json:"footer_active,omitempty"`
	FooterHTMLBody       string `json:"footer_html_body,omitempty"`
	FooterPlainBody      string `json:"footer_plain_body,omitempty"`
	LocalPart            string `json:"local_part,omitempty"`
	MayAccessImap        bool   `json:"may_access_imap,omitempty"`
	MayAccessManagesieve bool   `json:"may_access_managesieve,omitempty"`
	MayAccessPop3        bool   `json:"may_access_pop3,omitempty"`
	MayReceive           bool   `json:"may_receive,omitempty"`
	MaySend              bool   `json:"may_send,omitempty"`
	Name                 string `json:"name,omitempty"`
	Password             string `json:"password,omitempty"`
}

// Create returns a request-safe copy for identity create operations.
func (i Identity) Create() Identity {
	i.Address = ""
	i.DomainName = ""
	return i
}

// Update returns a request-safe copy for identity update operations.
func (i Identity) Update() Identity {
	i = i.Create()
	i.LocalPart = ""
	return i
}

// ListIdentities lists all the identities for the given domain and mailbox.
// It returns a slice of Identity structs and any error encountered.
func (c *Client) ListIdentities(ctx context.Context, domain *Domain, mailbox string) ([]Identity, error) {

	var identityList struct {
		Identities []Identity `json:"identities,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities", escapePathSegment(domain.Name), escapePathSegment(mailbox)))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &identityList); err != nil {
		return nil, err
	}

	return identityList.Identities, nil
}

// GetIdentity retrieves a single identity given the domain, mailbox and local part name.
// It returns a pointer to a Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, domain *Domain, mailbox string, localPart string) (*Identity, error) {

	var identity Identity

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(domain.Name), escapePathSegment(mailbox), escapePathSegment(localPart)))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &identity); err != nil {
		return nil, err
	}

	return &identity, nil
}

// NewIdentity creates a new identity.
// It returns a pointer to an Identity struct and any error encountered.
func (c *Client) NewIdentity(ctx context.Context, domain *Domain, mailbox string, identity *Identity) (*Identity, error) {

	jsonBody, err := json.Marshal(identity.Create())
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities", escapePathSegment(domain.Name), escapePathSegment(mailbox)), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, identity); err != nil {
		return nil, err
	}
	return identity, nil
}

// UpdateIdentity updates an identity in place given the domain, mailbox and a pointer to an Identity struct.
// It returns a pointer to a new Identity struct and any error encountered.
func (c *Client) UpdateIdentity(ctx context.Context, domain *Domain, mailbox string, i *Identity) (*Identity, error) {
	jsonBody, err := json.Marshal(i.Update())
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(domain.Name), escapePathSegment(mailbox), escapePathSegment(i.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, i); err != nil {
		return nil, err
	}
	return i, nil
}

// DeleteIdentity deletes an identity given the domain, mailbox and a pointer to an Identity struct.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, domain *Domain, mailbox string, i *Identity) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/mailboxes/%s/identities/%s", escapePathSegment(domain.Name), escapePathSegment(mailbox), escapePathSegment(i.LocalPart)))
	if err != nil {
		return err
	}
	return nil
}
