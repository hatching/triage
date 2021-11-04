// Copyright (C) 2020-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	_, err := client.CreateProfile(ctx, "pf1", []string{}, "notwork", 30)
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/profiles" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "POST" {
		t.Fatal("Unexpected method")
	}
}

func TestDeleteProfile(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if err := client.DeleteProfile(ctx, "pf1"); err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/profiles/pf1" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "DELETE" {
		t.Fatal("Unexpected method")
	}
}
