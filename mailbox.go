package migadu

import (
	"context"
	"encoding/json"
	"fmt"
)

// validateMailboxPasswordMethod validates password method requirements for mailboxes.
func validateMailboxPasswordMethod(mb *Mailbox, isCreate bool) error {
	if mb.PasswordMethod == "invitation" {
		if mb.PasswordRecoveryEmail == "" {
			return fmt.Errorf("password_recovery_email is required when password_method is 'invitation'")
		}
		if mb.Password != "" {
			return fmt.Errorf("password must not be provided when password_method is 'invitation'")
		}
	} else if isCreate && (mb.PasswordMethod == "" || mb.PasswordMethod == "password") {
		if mb.Password == "" {
			return fmt.Errorf("password is required when password_method is 'password' or not specified")
		}
	}
	return nil
}

// Mailbox represents a mailbox in the Migadu API.
//
// SpamAggressiveness is a string with the following valid values:
//
//	"strictest"       (most aggressive filtering)
//	"stricter"
//	"strict"
//	"default"         (use domain setting)
//	"permissive"
//	"more permissive"
//	"most permissive" (least aggressive filtering)
type Mailbox struct {
	ActivatedAt           string       `json:"activated_at,omitempty" api:"read-only"`
	Address               string       `json:"address,omitempty" api:"read-only"`
	ChangedAt             string       `json:"changed_at,omitempty" api:"read-only"`
	DailyIncomingLimit    int          `json:"daily_incoming_limit,omitempty"`
	DailyOutgoingLimit    int          `json:"daily_outgoing_limit,omitempty"`
	Delegations           []string     `json:"delegations,omitempty"`
	DomainName            string       `json:"domain_name,omitempty" api:"read-only"`
	Expireable            bool         `json:"expireable,omitempty"`
	ExpiresOn             string       `json:"expires_on,omitempty"`
	FooterActive          bool         `json:"footer_active,omitempty"`
	FooterHTMLBody        string       `json:"footer_html_body,omitempty"`
	FooterPlainBody       string       `json:"footer_plain_body,omitempty"`
	ForwardingTo          string       `json:"forwarding_to,omitempty" api:"create-only"`
	Forwardings           []Forwarding `json:"forwardings,omitempty"`
	Identities            []Identity   `json:"identities,omitempty"`
	IsActive              bool         `json:"is_active,omitempty" api:"read-only"`
	IsInternal            bool         `json:"is_internal,omitempty"`
	LastLoginAt           string       `json:"last_login_at,omitempty" api:"read-only"`
	LocalPart             string       `json:"local_part,omitempty" api:"create-only"`
	MayAccessImap         bool         `json:"may_access_imap,omitempty"`
	MayAccessManagesieve  bool         `json:"may_access_managesieve,omitempty"`
	MayAccessPop3         bool         `json:"may_access_pop3,omitempty"`
	MayReceive            bool         `json:"may_receive,omitempty"`
	MaySend               bool         `json:"may_send,omitempty"`
	MonthlyIncomingLimit  int          `json:"monthly_incoming_limit,omitempty"`
	MonthlyOutgoingLimit  int          `json:"monthly_outgoing_limit,omitempty"`
	Name                  string       `json:"name,omitempty"`
	Password              string       `json:"password,omitempty"`
	PasswordMethod        string       `json:"password_method,omitempty"`
	PasswordRecoveryEmail string       `json:"password_recovery_email,omitempty"`
	RecipientDenylist     []string     `json:"recipient_denylist,omitempty"`
	RemoveUponExpiry      bool         `json:"remove_upon_expiry,omitempty"`
	SenderAllowlist       []string     `json:"sender_allowlist,omitempty"`
	SenderDenylist        []string     `json:"sender_denylist,omitempty"`
	SpamAction            string       `json:"spam_action,omitempty"`
	SpamAggressiveness    string       `json:"spam_aggressiveness,omitempty"`
	StorageUsage          float64      `json:"storage_usage,omitempty" api:"read-only"`
	WeeklyIncomingLimit   int          `json:"weekly_incoming_limit,omitempty"`
	WeeklyOutgoingLimit   int          `json:"weekly_outgoing_limit,omitempty"`
	WildcardSender        bool         `json:"wildcard_sender,omitempty"`
}

// ListMailboxes lists all the mailboxes for the given domain.
// It returns a slice of Mailbox structs and any error encountered.
func (c *Client) ListMailboxes(ctx context.Context, d *Domain) ([]Mailbox, error) {
	var mailboxList struct {
		Mailboxes []Mailbox `json:"mailboxes,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes", d.Name))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &mailboxList); err != nil {
		return nil, err
	}

	return mailboxList.Mailboxes, nil
}

// GetMailbox retrieves a single mailbox given the domain and a Mailbox with LocalPart set.
// It returns a pointer to a Mailbox struct and any error encountered.
func (c *Client) GetMailbox(ctx context.Context, d *Domain, mb *Mailbox) (*Mailbox, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/mailboxes/%s", d.Name, escapePathSegment(mb.LocalPart)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, mb); err != nil {
		return nil, err
	}

	return mb, nil
}

// NewMailbox creates a new mailbox.
// It returns a pointer to a Mailbox struct and any error encountered.
func (c *Client) NewMailbox(ctx context.Context, d *Domain, mb *Mailbox) (*Mailbox, error) {
	if err := validateMailboxPasswordMethod(mb, true); err != nil {
		return nil, err
	}

	transformed, err := Transform(*mb, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/mailboxes", d.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, mb); err != nil {
		return nil, err
	}

	return mb, nil
}

// UpdateMailbox updates a mailbox in place given the domain and a pointer to a Mailbox struct.
// It returns a pointer to a new Mailbox struct and any error encountered.
func (c *Client) UpdateMailbox(ctx context.Context, d *Domain, mb *Mailbox) (*Mailbox, error) {
	if err := validateMailboxPasswordMethod(mb, false); err != nil {
		return nil, err
	}

	transformed, err := Transform(*mb, "update")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/mailboxes/%s", d.Name, escapePathSegment(mb.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, mb); err != nil {
		return nil, err
	}

	return mb, nil
}

// DeleteMailbox deletes a mailbox given the domain and a pointer to a Mailbox struct.
// It returns any error encountered.
func (c *Client) DeleteMailbox(ctx context.Context, d *Domain, mb *Mailbox) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/mailboxes/%s", d.Name, escapePathSegment(mb.LocalPart)))
	if err != nil {
		return err
	}
	return nil
}
