// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"testing"
)

func TestStaticReport(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	_, err := client.SampleStaticReport(ctx, "test-123")
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/reports/static" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestOverviewReport(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	_, err := client.SampleOverviewReport(ctx, "test-123")
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v1/samples/test-123/overview.json" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestTaskReport(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	_, err := client.SampleTaskReport(ctx, "test-123", "task-5")
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/task-5/report_triage.json" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestTaskFile(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	_, err := client.SampleTaskFile(ctx, "test-123", "task-5", "file-x")
	if err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/task-5/file-x" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestArchiveTar(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if _, err := client.SampleArchiveTAR(ctx, "test-123"); err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/archive" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestArchiveZIP(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	if _, err := client.SampleArchiveZIP(ctx, "test-123"); err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/archive.zip" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}

func TestSamplePath(t *testing.T) {
	ctx := context.Background()
	m := ClientMock{}
	m.Response = nil
	m.StatusCode = 200
	client := &Client{
		httpClient: &m,
	}
	type test struct {
		Name string
	}
	if err := client.SamplePath(ctx, "test-123", "/output.json", &test{}); err != nil {
		t.Fatal(err)
	}
	if m.RequestUrl != "/v0/samples/test-123/output.json" {
		t.Fatalf("Expected other request url %v", m.RequestUrl)
	}
	if m.RequestMethod != "GET" {
		t.Fatal("Unexpected method")
	}
}
