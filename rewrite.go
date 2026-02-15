package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Rewrite represents a rewrite rule in the Migadu API.
type Rewrite struct {
	Destinations  []string `json:"destinations,omitempty"`
	LocalPartRule string   `json:"local_part_rule,omitempty"`
	Name          string   `json:"name,omitempty"`
	OrderNum      int      `json:"order_num,omitempty"`
}

// rewriteJSON is used when a new/updated alias object to the API.
type rewriteJSON struct {
	Rewrite
	DestinationsJSON string `json:"destinations,omitempty"`
}

// convertDestinationsField takes a slice of strings and joins them into a comma separated line.
func (r *rewriteJSON) convertDestinationsField() {
	r.DestinationsJSON = strings.Join(r.Destinations, ",")
	r.Destinations = nil
}

// ListRewrites lists all the rewrites for the given domain.
// It returns a slice of Rewrite structs and any error encountered.
func (c *Client) ListRewrites(ctx context.Context, domain *Domain) ([]Rewrite, error) {

	var rewriteList struct {
		Rewrites []Rewrite `json:"rewrites,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites", domain.Name))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &rewriteList); err != nil {
		return nil, err
	}

	return rewriteList.Rewrites, nil
}

// GetRewrite retrieves a single rewrite given the domain and its name.
// It returns a pointer to an Rewrite struct and any error encountered.
func (c *Client) GetRewrite(ctx context.Context, domain *Domain, name string) (*Rewrite, error) {

	var rewrite Rewrite

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites/%s", domain.Name, name))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &rewrite); err != nil {
		return nil, err
	}

	return &rewrite, nil
}

// NewRewrite creates a new rewrite.
// It returns a pointer to an Rewrite struct and any error encountered.
func (c *Client) NewRewrite(ctx context.Context, domain *Domain, rewrite *Rewrite) (*Rewrite, error) {
	rewriteJSON := rewriteJSON{Rewrite: *rewrite}
	rewriteJSON.convertDestinationsField()
	jsonBody, err := json.Marshal(rewriteJSON)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/rewrites", domain.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, rewrite); err != nil {
		return nil, err
	}
	return rewrite, nil
}

// UpdateRewrite updates an rewrite in place given the domain and a pointer to an Rewrite struct.
// It returns a pointer to a new Rewrite struct and any error encountered.
func (c *Client) UpdateRewrite(ctx context.Context, domain *Domain, r *Rewrite) (*Rewrite, error) {
	rewriteJSON := rewriteJSON{Rewrite: *r}
	rewriteJSON.convertDestinationsField()
	jsonBody, err := json.Marshal(rewriteJSON)
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/rewrites/%s", domain.Name, r.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteRewrite deletes an rewrite given the domain and a pointer to an Rewrite struct.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, domain *Domain, r *Rewrite) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/rewrites/%s", domain.Name, r.Name))
	if err != nil {
		return err
	}
	return nil
}
