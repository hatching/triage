// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"net/http"
	"time"
)

type Profile struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Network string   `json:"network"`
	Timeout uint     `json:"timeout"`
}

func (c *Client) CreateProfile(ctx context.Context, name string, tags []string, network string, timeout time.Duration) (*Profile, error) {
	req := Profile{
		Name:    name,
		Tags:    tags,
		Network: network,
		Timeout: uint(timeout / time.Second),
	}
	var resp Profile
	err := c.jsonRequestJSON(ctx, http.MethodPost, "/v0/profiles", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteProfile deletes a profile by its ID or name.
func (c *Client) DeleteProfile(ctx context.Context, id string) error {
	return c.jsonRequestJSON(
		ctx, http.MethodDelete, "/v0/profiles/"+id, nil, nil,
	)
}

func (c *Client) Profiles(ctx context.Context) ([]Profile, error) {
	var response struct {
		Data []Profile `json:"data"`
	}
	err := c.jsonRequestJSON(
		ctx, http.MethodGet, "/v0/profiles", nil, &response,
	)
	if err != nil {
		return nil, err
	}
	return response.Data, err
}
