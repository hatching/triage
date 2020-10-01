// Copyright (C) 2020 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"testing"

	mock_triage "github.com/hatching/triage/go/mock"
)

func TestCreateProfile(t *testing.T) {
	ctx := context.Background()
	m := mock_triage.ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if _, err := client.CreateProfile(ctx, "pf1", []string{}, "notwork", 30); err != nil {
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
	m := mock_triage.ClientMock{}
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
