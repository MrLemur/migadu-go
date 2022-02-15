package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// ListIdentities lists all the identities for the given mailbox local part name.
// Ir returns a pointer to an array of Identity structs and any error encountered.
func (c *Client) ListIdentities(ctx context.Context, mailbox string) (*[]Identity, error) {

	var identityList struct {
		Indentities []Identity `json:"identities,omitempty,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities", mailbox))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &identityList)

	return &identityList.Indentities, nil
}

// GetIdentity  retrieves a single identity given its mailbox name and local part name.
// It returns a pointer to a Identity struct and any error encountered.
func (c *Client) GetIdentity(ctx context.Context, mailbox string, localPart string) (*Identity, error) {

	var identity Identity

	resp, err := c.Get(ctx, fmt.Sprintf("mailboxes/%s/identities/%s", mailbox, localPart))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &identity)

	return &identity, nil
}

// NewIdentity creates a new identity given the mailbox, local part name and a display name.
// It returns a pointer to am Identity struct and any error encountered.
func (c *Client) NewIdentity(ctx context.Context, mailbox string, localPart string, displayName string) (*Identity, error) {

	var identity = Identity{LocalPart: localPart, Name: displayName}

	jsonBody, _ := json.Marshal(identity)
	resp, err := c.Post(ctx, fmt.Sprintf("mailboxes/%s/identities", mailbox), jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &identity)
	return &identity, nil
}

// UpdateIdentity updates an identity in place given a pointer to an Identity struct.
// It returns a pointer to a new Identity struct and any error encountered.
func (c *Client) UpdateIdentity(ctx context.Context, mailbox string, i *Identity) (*Identity, error) {
	jsonBody, _ := json.Marshal(i)
	resp, err := c.Put(ctx, fmt.Sprintf("mailboxes/%s/identities/%s", mailbox, i.LocalPart), jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &i)
	return i, nil
}

// DeleteIdentity deletes an identity given a pointer to an Identity struct.
// It returns any error encountered.
func (c *Client) DeleteIdentity(ctx context.Context, mailbox string, i *Identity) error {
	_, err := c.Delete(ctx, fmt.Sprintf("mailboxes/%s/identities/%s", mailbox, i.LocalPart))
	if err != nil {
		return err
	}
	return nil
}
