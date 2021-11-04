// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"strings"
	"testing"
)

func TestSubmitUrl(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	t.Run("test success", func(t *testing.T) {
		_, err := client.SubmitSampleURL(
			ctx, "http://google.com", false, []ProfileSelection{},
		)
		if err != nil {
			t.Error(err)
		}
		if m.RequestMethod != "POST" {
			t.Fatal("Unexpected method")
		}
		if m.RequestUrl != "/v0/samples" {
			t.Fatal("Unexpected endpoint")
		}
		if m.RequestBody != `{"interactive":false,"kind":"url","profiles":[],"url":"http://google.com"}` {
			t.Fatal("Unexpected request body")
		}
	})
	t.Run("test fail", func(t *testing.T) {
		m.StatusCode = 400
		_, err := client.SubmitSampleURL(
			ctx, "http://google.com", false, []ProfileSelection{},
		)
		if !strings.Contains(err.Error(), "400") {
			t.Fatalf("expected error")
		}
	})
}

func TestSearchForUser(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = struct {
		Next string   `json:"next"`
		Data []Sample `json:"data"`
	}{
		Next: "2006-01-02T15:04:05Z07:00",
		Data: []Sample{
			{
				ID: "test1",
			},
			{
				ID: "test2",
			},
		},
	}
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	samples := client.Search(ctx, "NOT family:emotet", 2000)
	<-samples
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
	if m.RequestUrl != "/v0/search?query=NOT+family%3Aemotet&limit=200" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	<-samples
	<-samples
	if m.RequestUrl != "/v0/search?query=NOT+family%3Aemotet&limit=200&offset=2006-01-02T15:04:05Z07:00" {
		t.Fatalf("Expected other request url")
	}
}

func TestSamplesForUser(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = struct {
		Next string   `json:"next"`
		Data []Sample `json:"data"`
	}{
		Next: "2006-01-02T15:04:05Z07:00",
		Data: []Sample{
			{
				ID: "test1",
			},
			{
				ID: "test2",
			},
		},
	}
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	samples := client.SamplesForUser(ctx, 2000)
	<-samples
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
	if m.RequestUrl != "/v0/samples?subset=owned&limit=200" {
		t.Fatalf("Expected other request url")
	}
	<-samples
	<-samples
	if m.RequestUrl != "/v0/samples?subset=owned&limit=200&offset=2006-01-02T15:04:05Z07:00" {
		t.Fatalf("Expected other request url")
	}
}

func TestPublicSamples(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = struct {
		Next string   `json:"next"`
		Data []Sample `json:"data"`
	}{
		Next: "2006-01-02T15:04:05Z07:00",
		Data: []Sample{{
			ID: "test1",
		}, {
			ID: "test2",
		}},
	}
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	samples := client.PublicSamples(ctx, 2000)
	<-samples
	if m.RequestUrl != "/v0/samples?subset=public&limit=200" {
		t.Fatalf("Expected other request url")
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
	<-samples
	<-samples
	if m.RequestUrl != "/v0/samples?subset=public&limit=200&offset=2006-01-02T15:04:05Z07:00" {
		t.Fatalf("Expected other request url")
	}
}

func TestSampleById(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if _, err := client.SampleByID(ctx, "test-123"); err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
	m.StatusCode = 400
	_, err := client.SampleByID(ctx, "test-123")
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("expected error")
	}
}

func TestSetSampleProfile(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	err := client.SetSampleProfile(ctx, "test-123", []ProfileSelection{})
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestMethod != "POST" {
		t.Fatal("Unexpected method")
	}
	if m.RequestBody != `{"auto":false,"profiles":[]}` {
		t.Fatal("Unexpected request body")
	}
	if m.RequestUrl != "/v0/samples/test-123/profile" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
}

func TestSetSampleProfileAutomatically(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	err := client.SetSampleProfileAutomatically(
		ctx, "test-123", []string{"file1.exe"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestMethod != "POST" {
		t.Fatal("Unexpected method")
	}
	if m.RequestBody != `{"auto":true,"pick":["file1.exe"]}` {
		t.Fatal("Unexpected request body")
	}
	if m.RequestUrl != "/v0/samples/test-123/profile" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
}

func TestDeleteSample(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if err := client.DeleteSample(ctx, "test-123"); err != nil {
		t.Fatal(err)
	}
	if m.RequestMethod != "DELETE" {
		t.Fatal("Unexpected method")
	}
	if m.RequestUrl != "/v0/samples/test-123" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
}
