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
	Address      string `json:"address,omitempty"`
	Destinations string `json:"destinations,omitempty"`
	DomainName   string `json:"domain_name,omitempty"`
	IsInternal   bool   `json:"is_internal,omitempty"`
	LocalPart    string `json:"local_part,omitempty"`
}

func (a *Alias) UnmarshalJSON(data []byte) error {
	type alias Alias
	type rawAlias struct {
		alias
		Destinations json.RawMessage `json:"destinations"`
	}

	var raw rawAlias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*a = Alias(raw.alias)
	dest, err := parseDestinations(raw.Destinations)
	if err != nil {
		return err
	}
	a.Destinations = dest
	return nil
}

func parseDestinations(data json.RawMessage) (string, error) {
	if len(data) == 0 || string(data) == "null" {
		return "", nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return s, nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		return strings.Join(arr, ","), nil
	}

	return "", fmt.Errorf("unsupported destinations format: %s", string(data))
}

// Create returns a request-safe copy for alias create operations.
func (a Alias) Create() Alias {
	a.Address = ""
	a.DomainName = ""
	return a
}

// Update returns a request-safe copy for alias update operations.
func (a Alias) Update() Alias {
	a = a.Create()
	a.LocalPart = ""
	return a
}

// ListAliases lists all the aliases for the given domain.
// It returns a slice of Alias structs and any error encountered.
func (c *Client) ListAliases(ctx context.Context, domain *Domain) ([]Alias, error) {

	var aliasList struct {
		Aliases []Alias `json:"address_aliases,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases", escapePathSegment(domain.Name)))
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

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(domain.Name), escapePathSegment(localPart)))
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
	jsonBody, err := json.Marshal(alias.Create())
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/aliases", escapePathSegment(domain.Name)), jsonBody)
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
	jsonBody, err := json.Marshal(a.Update())
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(domain.Name), escapePathSegment(a.LocalPart)), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteAlias deletes an alias given the domain and a pointer to an Alias struct.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, domain *Domain, a *Alias) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/aliases/%s", escapePathSegment(domain.Name), escapePathSegment(a.LocalPart)))
	if err != nil {
		return err
	}
	return nil
}
