package migadu

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// convertDestinationsField takes a slice of strings and joins them into a comma seperated line.
func (r *rewriteJSON) convertDestinationsField() {
	r.DestinationsJSON = strings.Join(r.Destinations, ",")
	r.Destinations = nil
}

// ListRewrites lists all the rewrites for the domain configured on the client.
// Ir returns a pointer to an array of Rewrite structs and any error encountered.
func (c *Client) ListRewrites(ctx context.Context) (*[]Rewrite, error) {

	var rewriteList struct {
		Rewrites []Rewrite `json:"rewrites,omitempty"`
	}

	resp, err := c.Get(ctx, "rewrites")
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &rewriteList)

	return &rewriteList.Rewrites, nil
}

// GetRewrite retrieves a single rewrite given its name.
// It returns a pointer to an Rewrite struct and any error encountered.
func (c *Client) GetRewrite(ctx context.Context, name string) (*Rewrite, error) {

	var rewrite Rewrite

	resp, err := c.Get(ctx, fmt.Sprintf("rewrites/%s", name))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &rewrite)

	return &rewrite, nil
}

// NewRewrite creates a new rewrite given the name, local part rule and its destinations.
// It returns a pointer to an Rewrite struct and any error encountered.
func (c *Client) NewRewrite(ctx context.Context, name string, localPartRule string, destinations []string) (*Rewrite, error) {
	var rewrite = Rewrite{Name: name, LocalPartRule: localPartRule, Destinations: destinations}
	rewriteJSON := rewriteJSON{Rewrite: rewrite}
	rewriteJSON.convertDestinationsField()
	jsonBody, _ := json.Marshal(rewriteJSON)
	resp, err := c.Post(ctx, "rewrites", jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &rewrite)
	return &rewrite, nil
}

// UpdateRewrite updates an rewrite in place given a pointer to an Rewrite struct.
// It returns a pointer to a new Rewrite struct and any error encountered.
func (c *Client) UpdateRewrite(ctx context.Context, r *Rewrite) (*Rewrite, error) {
	rewriteJSON := rewriteJSON{Rewrite: *r}
	rewriteJSON.convertDestinationsField()
	jsonBody, _ := json.Marshal(rewriteJSON)
	resp, err := c.Put(ctx, fmt.Sprintf("rewrites/%s", r.Name), jsonBody)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &r)
	return r, nil
}

// DeleteRewrite deletes an rewrite given a pointer to an Rewrite struct.
// It returns any error encountered.
func (c *Client) DeleteRewrite(ctx context.Context, r *Rewrite) error {
	_, err := c.Delete(ctx, fmt.Sprintf("rewrites/%s", r.Name))
	if err != nil {
		return err
	}
	return nil
}
