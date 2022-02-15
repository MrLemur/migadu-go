package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// convertDestinationsField takes a slice of strings and joins them into a comma seperated line.
func (a *aliasJSON) convertDestinationsField() {
	a.DestinationsJSON = strings.Join(a.Destinations, ",")
	a.Destinations = nil
}

// ListAliases lists all the aliases for the domain configured on the client.
// Ir returns a pointer to an array of Alias structs and any error encountered.
func (c *Client) ListAliases(ctx context.Context) (*[]Alias, error) {

	var aliasList struct {
		Aliases []Alias `json:"address_aliases,omitempty"`
	}

	resp, err := c.Get(ctx, "aliases")
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &aliasList)

	return &aliasList.Aliases, nil
}

// GetAlias retrieves a single alias given its local part name.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) GetAlias(ctx context.Context, localPart string) (*Alias, error) {

	var alias Alias

	resp, err := c.Get(ctx, fmt.Sprintf("aliases/%s", localPart))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &alias)

	return &alias, nil
}

// NewAlias creates a new alias given the local part name and its destinations.
// It returns a pointer to an Alias struct and any error encountered.
func (c *Client) NewAlias(ctx context.Context, localPart string, destinations []string) (*Alias, error) {
	var alias = Alias{LocalPart: localPart, Destinations: destinations}
	aliasJSON := aliasJSON{Alias: alias}
	aliasJSON.convertDestinationsField()
	jsonBody, _ := json.Marshal(aliasJSON)
	resp, err := c.Post(ctx, "aliases", jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &alias)
	return &alias, nil
}

// UpdateAlias updates an alias in place given a pointer to an Alias struct.
// It returns a pointer to a new Alias struct and any error encountered.
func (c *Client) UpdateAlias(ctx context.Context, a *Alias) (*Alias, error) {
	aliasJSON := aliasJSON{Alias: *a}
	aliasJSON.convertDestinationsField()
	jsonBody, _ := json.Marshal(aliasJSON)
	resp, err := c.Put(ctx, fmt.Sprintf("aliases/%s", a.LocalPart), jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &a)
	return a, nil
}

// DeleteAlias deletes an alias given a pointer to an Alias struct.
// It returns any error encountered.
func (c *Client) DeleteAlias(ctx context.Context, a *Alias) error {
	_, err := c.Delete(ctx, fmt.Sprintf("aliases/%s", a.LocalPart))
	if err != nil {
		return err
	}
	return nil
}
