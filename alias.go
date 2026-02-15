package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Alias represents an alias in the Migadu API.
type Alias struct {
	Address          string   `json:"address,omitempty"`
	Destinations     []string `json:"destinations,omitempty"`
	DomainName       string   `json:"domain_name,omitempty"`
	Expireable       bool     `json:"expireable,omitempty"`
	ExpiresOn        string   `json:"expires_on,omitempty"`
	IsInternal       bool     `json:"is_internal,omitempty"`
	LocalPart        string   `json:"local_part,omitempty"`
	RemoveUponExpiry bool     `json:"remove_upon_expiry,omitempty"`
}

// aliasJSON is used when a new/updated alias object to the API.
type aliasJSON struct {
	Alias
	DestinationsJSON string `json:"destinations,omitempty"`
}

// convertDestinationsField takes a slice of strings and joins them into a comma separated line.
func (a *aliasJSON) convertDestinationsField() {
	a.DestinationsJSON = strings.Join(a.Destinations, ",")
	a.Destinations = nil
}

// ListAliases lists all the aliases for the given domain.
// It returns a slice of Alias structs and any error encountered.
func (c *Client) ListAliases(ctx context.Context, domain *Domain) ([]Alias, error) {

	var aliasList struct {
		Aliases []Alias `json:"address_aliases,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases", domain.Name))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &aliasList); err != nil {
		return nil, err
	}

	return aliasList.Aliases, nil
}

// GetAlias retrieves a single alias given the domain and local part name.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) GetAlias(ctx context.Context, domain *Domain, localPart string) (*Alias, error) {

	var alias Alias

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases/%s", domain.Name, localPart))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &alias); err != nil {
		return nil, err
	}

	return &alias, nil
}

// NewAlias creates a new alias.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) NewAlias(ctx context.Context, domain *Domain, alias *Alias) (*Alias, error) {
	aliasJSON := aliasJSON{Alias: *alias}
	aliasJSON.convertDestinationsField()
	jsonBody, err := json.Marshal(aliasJSON)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/aliases", domain.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, alias); err != nil {
		return nil, err
	}
	return alias, nil
}

// UpdateAlias updates an alias in place given the domain and a pointer to an Alias struct.
// It returns a pointer to a new Alias struct and any error encountered.
func (c *Client) UpdateAlias(ctx context.Context, domain *Domain, a *Alias) (*Alias, error) {
	aliasJSON := aliasJSON{Alias: *a}
	aliasJSON.convertDestinationsField()
	jsonBody, err := json.Marshal(aliasJSON)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/aliases/%s", domain.Name, a.LocalPart), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAlias deletes an alias given the domain and a pointer to an Alias struct.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, domain *Domain, a *Alias) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/aliases/%s", domain.Name, a.LocalPart))
	if err != nil {
		return err
	}
	return nil
}
