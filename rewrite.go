package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// validateRewriteDestinations validates that rewrite destinations are on the same domain.
func validateRewriteDestinations(r *Rewrite, d *Domain) error {
	if len(r.Destinations) == 0 {
		return fmt.Errorf("rewrite must have at least one destination")
	}

	for _, dest := range r.Destinations {
		if !strings.Contains(dest, "@") {
			return fmt.Errorf("invalid destination format: %s", dest)
		}

		parts := strings.Split(dest, "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid destination format: %s", dest)
		}

		if parts[1] != d.Name {
			return fmt.Errorf("rewrite destinations must be on the same domain (%s), got: %s (external domains will be rewritten by the API)", d.Name, dest)
		}
	}
	return nil
}

// Rewrite represents a rewrite rule in the Migadu API.
type Rewrite struct {
	Destinations  []string `json:"destinations,omitempty"`
	LocalPartRule string   `json:"local_part_rule,omitempty"`
	Name          string   `json:"name,omitempty"`
	OrderNum      int      `json:"order_num,omitempty"`
}

// ListRewrites lists all the rewrites for the given domain.
// It returns a slice of Rewrite structs and any error encountered.
func (c *Client) ListRewrites(ctx context.Context, d *Domain) ([]Rewrite, error) {
	var rewriteList struct {
		Rewrites []Rewrite `json:"rewrites,omitempty"`
	}

	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites", escapePathSegment(d.Name)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, &rewriteList); err != nil {
		return nil, err
	}

	return rewriteList.Rewrites, nil
}

// GetRewrite retrieves a single rewrite given the domain and a Rewrite with Name set.
// It returns a pointer to a Rewrite struct and any error encountered.
func (c *Client) GetRewrite(ctx context.Context, d *Domain, r *Rewrite) (*Rewrite, error) {
	resp, err := c.Get(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(d.Name), escapePathSegment(r.Name)))
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, r); err != nil {
		return nil, err
	}

	return r, nil
}

// NewRewrite creates a new rewrite.
// It returns a pointer to a Rewrite struct and any error encountered.
func (c *Client) NewRewrite(ctx context.Context, d *Domain, r *Rewrite) (*Rewrite, error) {
	if err := validateRewriteDestinations(r, d); err != nil {
		return nil, err
	}

	transformed, err := Transform(*r, "create")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, fmt.Sprintf("domains/%s/rewrites", escapePathSegment(d.Name)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, r); err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateRewrite updates a rewrite in place given the domain and a pointer to a Rewrite struct.
// It returns a pointer to a new Rewrite struct and any error encountered.
func (c *Client) UpdateRewrite(ctx context.Context, d *Domain, r *Rewrite) (*Rewrite, error) {
	if err := validateRewriteDestinations(r, d); err != nil {
		return nil, err
	}

	transformed, err := Transform(*r, "update")
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(transformed)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(d.Name), escapePathSegment(r.Name)), jsonBody)
	if err != nil {
		return nil, err
	}
	if err = readAndUnmarshal(resp.Body, r); err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteRewrite deletes a rewrite given the domain and a pointer to a Rewrite struct.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, d *Domain, r *Rewrite) error {
	_, err := c.Delete(ctx, fmt.Sprintf("domains/%s/rewrites/%s", escapePathSegment(d.Name), escapePathSegment(r.Name)))
	if err != nil {
		return err
	}
	return nil
}
