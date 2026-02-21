package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// validateAliasDestinations validates that alias destinations are on the same domain.
func validateAliasDestinations(a *Alias, d *Domain) error {
	if len(a.Destinations) == 0 {
		return fmt.Errorf("alias must have at least one destination")
	}

	for _, dest := range a.Destinations {
		if !strings.Contains(dest, "@") {
			return fmt.Errorf("invalid destination format: %s", dest)
		}

		parts := strings.Split(dest, "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid destination format: %s", dest)
		}

		if parts[1] != d.Name {
			return fmt.Errorf("alias destinations must be on the same domain (%s), got: %s", d.Name, dest)
		}
	}
	return nil
}

// Alias represents an alias in the Migadu API.
type Alias struct {
	Address          string   `json:"address,omitempty" api:"read-only"`
	ExpiresOn        string   `json:"expires_on,omitempty"`
	Expireable       bool     `json:"expirable,omitempty"`
	Destinations     []string `json:"destinations,omitempty"`
	DomainName       string   `json:"domain_name,omitempty" api:"read-only"`
	IsInternal       bool     `json:"is_internal,omitempty"`
	RemoveUponExpiry bool     `json:"remove_upon_expiry,omitempty"`
	LocalPart        string   `json:"local_part,omitempty" api:"create-only"`
}

// AliasList is a wrapper for a list of aliases returned by the API.
type AliasList struct {
	Aliases []Alias `json:"address_aliases,omitempty"`
}

// ListAliases lists all the aliases for the given domain.
// It returns a slice of Alias structs and any error encountered.
func (c *Client) ListAliases(ctx context.Context, d *Domain) ([]Alias, error) {

	if d == nil {
		return nil, fmt.Errorf("domain must not be nil")
	}

	var aliasList AliasList

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases", escapePathSegment(d.Name)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &aliasList); err != nil {
		return nil, err
	}

	return aliasList.Aliases, nil
}

// GetAlias retrieves a single alias given the domain and local part name.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) GetAlias(ctx context.Context, d *Domain, a *Alias) (*Alias, error) {

	if d == nil {
		return nil, fmt.Errorf("domain must not be nil")
	}
	if a == nil {
		return nil, fmt.Errorf("alias must not be nil")
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(d.Name), escapePathSegment(a.LocalPart)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, a); err != nil {
		return nil, err
	}

	return a, nil
}

// NewAlias creates a new alias.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) NewAlias(ctx context.Context, d *Domain, a *Alias) (*Alias, error) {
	if d == nil {
		return nil, fmt.Errorf("domain must not be nil")
	}
	if a == nil {
		return nil, fmt.Errorf("alias must not be nil")
	}

	if err := validateAliasDestinations(a, d); err != nil {
		return nil, err
	}

	transformed, err := Transform(*a, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/aliases", escapePathSegment(d.Name)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, a); err != nil {
		return nil, err
	}

	return a, nil
}

// UpdateAlias updates an alias in place given the domain and a pointer to an Alias struct.
// It returns a pointer to a new Alias struct and any error encountered.
func (c *Client) UpdateAlias(ctx context.Context, d *Domain, a *Alias) (*Alias, error) {
	if d == nil {
		return nil, fmt.Errorf("domain must not be nil")
	}
	if a == nil {
		return nil, fmt.Errorf("alias must not be nil")
	}

	if err := validateAliasDestinations(a, d); err != nil {
		return nil, err
	}

	transformed, err := Transform(*a, "update")
	if err != nil {
		return nil, err
	}
	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(d.Name), escapePathSegment(a.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAlias deletes an alias given the domain and a pointer to an Alias struct.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, d *Domain, a *Alias) error {
	if d == nil {
		return fmt.Errorf("domain must not be nil")
	}
	if a == nil {
		return fmt.Errorf("alias must not be nil")
	}
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(d.Name), escapePathSegment(a.LocalPart)))
	if err != nil {
		return err
	}
	return nil
}
