package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Rewrite represents a rewrite rule in the Migadu API.
type Rewrite struct {
	Destinations  string `json:"destinations,omitempty"`
	LocalPartRule string `json:"local_part_rule,omitempty"`
	Name          string `json:"name,omitempty"`
	OrderNum      int    `json:"order_num,omitempty"`
}

func (r *Rewrite) UnmarshalJSON(data []byte) error {
	type alias Rewrite
	type rawRewrite struct {
		alias
		Destinations json.RawMessage `json:"destinations"`
	}

	var raw rawRewrite
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = Rewrite(raw.alias)
	dest, err := parseDestinations(raw.Destinations)
	if err != nil {
		return err
	}
	r.Destinations = dest
	return nil
}

// Create returns a request-safe copy for rewrite create operations.
func (r Rewrite) Create() Rewrite {
	// No read-only fields to strip for create
	return r
}

// Update returns a request-safe copy for rewrite update operations.
func (r Rewrite) Update() Rewrite {
	r = r.Create()
	r.Name = ""
	return r
}

// ListRewrites lists all the rewrites for the given domain.
// It returns a slice of Rewrite structs and any error encountered.
func (c *Client) ListRewrites(ctx context.Context, domain *Domain) ([]Rewrite, error) {

	var rewriteList struct {
		Rewrites []Rewrite `json:"rewrites,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites", escapePathSegment(domain.Name)))
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

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(domain.Name), escapePathSegment(name)))
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
	jsonBody, err := json.Marshal(rewrite.Create())
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/rewrites", escapePathSegment(domain.Name)), jsonBody)
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
	jsonBody, err := json.Marshal(r.Update())
	if err != nil {
		return nil, err
	}
	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(domain.Name), escapePathSegment(r.Name)), jsonBody)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, r); err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteRewrite deletes an rewrite given the domain and a pointer to an Rewrite struct.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, domain *Domain, r *Rewrite) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(domain.Name), escapePathSegment(r.Name)))
	if err != nil {
		return err
	}
	return nil
}
